package v1

import (
	"ToDoVerba/internal/dto"
	"ToDoVerba/internal/service"
	mockservice "ToDoVerba/internal/service/mocks"
	"ToDoVerba/pkg/logging"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHandler_taskCreate(t *testing.T) {
	type mockBehaviour func(s *mockservice.MockITaskService, user *dto.TaskCreate)

	parseTime := func(s string) pgtype.Timestamptz {
		t, _ := time.Parse(time.RFC3339, s)
		return pgtype.Timestamptz{
			Time:  t,
			Valid: true,
		}
	}

	testTable := []struct {
		name            string
		inputContType   string
		inputBody       string
		inputDTO        *dto.TaskCreate
		mockBehaviour   mockBehaviour
		expectedCode    int
		expectedBody    string
		bodyMustContain string
	}{
		{
			name: "201_valid_input",
			inputBody: `{
							"title": "First Task",
							"description": "First description",
							"due_date": "2024-09-05T15:04:05+05:00"
						}`,
			inputDTO: &dto.TaskCreate{
				Title:       "First Task",
				Description: "First description",
				DueDate:     parseTime("2024-09-05T15:04:05+05:00"),
			},
			inputContType: "application/json",
			mockBehaviour: func(s *mockservice.MockITaskService, user *dto.TaskCreate) {
				s.EXPECT().Create(user).Return(&dto.TaskRead{
					Id:          7,
					Title:       "First Task",
					Description: "First description",
					DueDate:     parseTime("2024-09-05T15:04:05+05:00"),
					CreatedAt:   parseTime("2022-09-05T15:04:05+05:00"),
					UpdatedAt:   parseTime("2023-09-05T15:04:05+05:00"),
				},
					nil,
				)
			},
			expectedCode: 201,
			expectedBody: `{
								"id": 7,
								"title": "First Task",
								"description": "First description",
								"due_date": "2024-09-05T15:04:05+05:00",
								"created_at": "2022-09-05T15:04:05+05:00",
								"updated_at": "2023-09-05T15:04:05+05:00"
							}`,
		},
		{
			name: "400_invalid_content_type",
			inputBody: `{
							"title": "First Task",
							"description": "First description",
							"due_date": "2024-09-05T15:04:05+05:00"
						}`,
			inputContType: "plain/text",
			inputDTO:      &dto.TaskCreate{},
			mockBehaviour: func(s *mockservice.MockITaskService, user *dto.TaskCreate) {
			},
			expectedCode: 400,
			expectedBody: `{"error":"content-type is not application/json"}`,
		},
		{
			name: "400_invalid_json_struct_input",
			inputBody: `{
							"title": "First Task,
							"description": "First description",
							"due_date": "2024-09-05T15:04:05+05:00"
						}`,
			inputContType: "application/json",
			inputDTO:      &dto.TaskCreate{},
			mockBehaviour: func(s *mockservice.MockITaskService, user *dto.TaskCreate) {},
			expectedCode:  400,
			expectedBody:  `{"error":"invalid character '\\n' in string literal"}`,
		},
		{
			name: "400_invalid_all_val_input",
			inputBody: `{
							"title": "",
							"description": "",
							"due_date": "2024-09-05T15:74:05+05:00"
						}`,
			inputContType: "application/json",
			inputDTO:      &dto.TaskCreate{},
			mockBehaviour: func(s *mockservice.MockITaskService, user *dto.TaskCreate) {},
			expectedCode:  400,
			expectedBody:  `{"error":"Title is required;Description is required;DueDate is required and must be in RFC3339 format;"}`,
		},
		{
			name: "400_invalid_due_date_input",
			inputBody: `{
							"title": "First Task",
							"description": "First description",
							"due_date": "2024-09-05T15:74:05+05:00"
						}`,
			inputContType: "application/json",
			inputDTO:      &dto.TaskCreate{},
			mockBehaviour: func(s *mockservice.MockITaskService, user *dto.TaskCreate) {},
			expectedCode:  400,
			expectedBody:  `{"error":"DueDate is required and must be in RFC3339 format;"}`,
		},
		{
			name: "500_unknown_error",
			inputBody: `{
							"title": "First Task",
							"description": "First description",
							"due_date": "2024-09-05T15:04:05+05:00"
						}`,
			inputContType: "application/json",
			inputDTO: &dto.TaskCreate{
				Title:       "First Task",
				Description: "First description",
				DueDate:     parseTime("2024-09-05T15:04:05+05:00"),
			},
			mockBehaviour: func(s *mockservice.MockITaskService, user *dto.TaskCreate) {
				s.EXPECT().Create(user).Return(&dto.TaskRead{}, errors.New("some error"))
			},
			expectedCode: 500,
			expectedBody: `{"error":"some error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()
			taskService := mockservice.NewMockITaskService(c)
			testCase.mockBehaviour(taskService, testCase.inputDTO)

			services := service.Services{Task: taskService}
			handler := NewHandler(Deps{
				Service: services,
				Logger:  logging.GetLoggerTest(),
			})

			//Test server
			r := httprouter.New()
			r.POST("/tasks", handler.taskCreate)

			//http test
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/tasks", strings.NewReader(testCase.inputBody))
			req.Header.Set("Content-Type", testCase.inputContType)

			//Perform request
			r.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, testCase.expectedCode, w.Code)
			if testCase.expectedBody != "" {
				if testCase.expectedBody == "nil" {
					testCase.expectedBody = ""
				}
				assert.JSONEq(t, testCase.expectedBody, w.Body.String())
			}
			if testCase.bodyMustContain != "" {
				assert.Contains(t, w.Body.String(), testCase.bodyMustContain)
			}
		})
	}
}

func TestHandler_taskList(t *testing.T) {
	type mockBehaviour func(s *mockservice.MockITaskService)

	parseTime := func(s string) pgtype.Timestamptz {
		t, _ := time.Parse(time.RFC3339, s)
		return pgtype.Timestamptz{
			Time:  t,
			Valid: true,
		}
	}

	testTable := []struct {
		name            string
		mockBehaviour   mockBehaviour
		expectedCode    int
		expectedBody    string
		bodyMustContain string
	}{
		{
			name: "200_multiple_tasks_response",
			mockBehaviour: func(s *mockservice.MockITaskService) {
				s.EXPECT().List().Return([]dto.TaskRead{
					{
						Id:          9,
						Title:       "First task",
						Description: "First description",
						DueDate:     parseTime("2024-09-05T15:04:05+05:00"),
						CreatedAt:   parseTime("2022-09-05T15:04:05+05:00"),
						UpdatedAt:   parseTime("2023-09-05T15:04:05+05:00"),
					},
					{
						Id:          10,
						Title:       "Second task",
						Description: "Second description",
						DueDate:     parseTime("2024-09-15T15:04:05+05:00"),
						CreatedAt:   parseTime("2022-09-15T15:04:05+05:00"),
						UpdatedAt:   parseTime("2023-09-15T15:04:05+05:00"),
					},
				},
					nil,
				)
			},
			expectedCode: 200,
			expectedBody: `[{
								"id": 9,
								"title": "First task",
								"description": "First description",
								"due_date": "2024-09-05T15:04:05+05:00",
								"created_at": "2022-09-05T15:04:05+05:00",
								"updated_at": "2023-09-05T15:04:05+05:00"
							},
							{
								"id": 10,
								"title": "Second task",
								"description": "Second description",
								"due_date": "2024-09-15T15:04:05+05:00",
								"created_at": "2022-09-15T15:04:05+05:00",
								"updated_at": "2023-09-15T15:04:05+05:00"
							}]`,
		},
		{
			name: "200_no_tasks_response",
			mockBehaviour: func(s *mockservice.MockITaskService) {
				s.EXPECT().List().Return([]dto.TaskRead{}, nil)
			},
			expectedCode: 200,
			expectedBody: `[]`,
		},
		{
			name: "500_unknown_error",
			mockBehaviour: func(s *mockservice.MockITaskService) {
				s.EXPECT().List().Return(nil, errors.New("some error"))
			},
			expectedCode: 500,
			expectedBody: `{"error":"some error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()
			taskService := mockservice.NewMockITaskService(c)
			testCase.mockBehaviour(taskService)

			services := service.Services{Task: taskService}
			handler := NewHandler(Deps{
				Service: services,
				Logger:  logging.GetLoggerTest(),
			})

			//Test server
			r := httprouter.New()
			r.GET("/tasks", handler.taskList)

			//http test
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/tasks", strings.NewReader(""))

			//Perform request
			r.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, testCase.expectedCode, w.Code)
			if testCase.expectedBody != "" {
				if testCase.expectedBody == "nil" {
					testCase.expectedBody = ""
				}
				assert.JSONEq(t, testCase.expectedBody, w.Body.String())
			}
			if testCase.bodyMustContain != "" {
				assert.Contains(t, w.Body.String(), testCase.bodyMustContain)
			}
		})
	}
}

func TestHandler_taskFindById(t *testing.T) {
	type mockBehaviour func(s *mockservice.MockITaskService, id int)

	parseTime := func(s string) pgtype.Timestamptz {
		t, _ := time.Parse(time.RFC3339, s)
		return pgtype.Timestamptz{
			Time:  t,
			Valid: true,
		}
	}

	testTable := []struct {
		name            string
		inputParam      string
		inputId         int
		mockBehaviour   mockBehaviour
		expectedCode    int
		expectedBody    string
		bodyMustContain string
	}{
		{
			name:       "200_valid_param",
			inputParam: "129",
			inputId:    129,
			mockBehaviour: func(s *mockservice.MockITaskService, id int) {
				s.EXPECT().FindByID(id).Return(&dto.TaskRead{
					Id:          129,
					Title:       "First Task",
					Description: "First description",
					DueDate:     parseTime("2024-09-05T15:04:05+05:00"),
					CreatedAt:   parseTime("2022-09-05T15:04:05+05:00"),
					UpdatedAt:   parseTime("2023-09-05T15:04:05+05:00"),
				},
					nil,
				)
			},
			expectedCode: 200,
			expectedBody: `{
								"id": 129,
								"title": "First Task",
								"description": "First description",
								"due_date": "2024-09-05T15:04:05+05:00",
								"created_at": "2022-09-05T15:04:05+05:00",
								"updated_at": "2023-09-05T15:04:05+05:00"
							}`,
		},
		{
			name:       "400_invalid_param",
			inputParam: "129f",
			inputId:    129,
			mockBehaviour: func(s *mockservice.MockITaskService, id int) {
			},
			expectedCode: 400,
			expectedBody: `{"error":"strconv.Atoi: parsing \"129f\": invalid syntax"}`,
		},
		{
			name:       "404_no_tasks_found",
			inputParam: "129",
			inputId:    129,
			mockBehaviour: func(s *mockservice.MockITaskService, id int) {
				s.EXPECT().FindByID(id).Return(&dto.TaskRead{}, pgx.ErrNoRows)
			},
			expectedCode: 404,
			expectedBody: `{"error":"no rows in result set"}`,
		},
		{
			name:       "500_unknown_error",
			inputParam: "129",
			inputId:    129,
			mockBehaviour: func(s *mockservice.MockITaskService, id int) {
				s.EXPECT().FindByID(id).Return(&dto.TaskRead{}, errors.New("some error"))
			},
			expectedCode: 500,
			expectedBody: `{"error":"some error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()
			taskService := mockservice.NewMockITaskService(c)
			testCase.mockBehaviour(taskService, testCase.inputId)

			services := service.Services{Task: taskService}
			handler := NewHandler(Deps{
				Service: services,
				Logger:  logging.GetLoggerTest(),
			})

			//Test server
			r := httprouter.New()
			r.GET("/tasks/:id", handler.taskFindById)

			//http test
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/tasks/"+testCase.inputParam, strings.NewReader(""))

			//Perform request
			r.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, testCase.expectedCode, w.Code)
			if testCase.expectedBody != "" {
				if testCase.expectedBody == "nil" {
					testCase.expectedBody = ""
				}
				assert.JSONEq(t, testCase.expectedBody, w.Body.String())
			}
			if testCase.bodyMustContain != "" {
				assert.Contains(t, w.Body.String(), testCase.bodyMustContain)
			}
		})
	}
}

func TestHandler_taskUpdateById(t *testing.T) {
	type mockBehaviour func(s *mockservice.MockITaskService, id int, update *dto.TaskUpdate)

	parseTime := func(s string) pgtype.Timestamptz {
		t, _ := time.Parse(time.RFC3339, s)
		return pgtype.Timestamptz{
			Time:  t,
			Valid: true,
		}
	}

	testTable := []struct {
		name            string
		inputParam      string
		inputBody       string
		inputContType   string
		inputId         int
		inputDTO        *dto.TaskUpdate
		mockBehaviour   mockBehaviour
		expectedCode    int
		expectedBody    string
		bodyMustContain string
	}{
		{
			name:       "200_valid_input",
			inputParam: "129",
			inputId:    129,
			inputBody: `{
							"title": "First Task",
							"description": "First description",
							"due_date": "2024-09-05T15:04:05+05:00"
						}`,
			inputDTO: &dto.TaskUpdate{
				Title:       "First Task",
				Description: "First description",
				DueDate:     parseTime("2024-09-05T15:04:05+05:00"),
			},
			inputContType: "application/json",
			mockBehaviour: func(s *mockservice.MockITaskService, id int, update *dto.TaskUpdate) {
				s.EXPECT().UpdateById(id, update).Return(&dto.TaskRead{
					Id:          129,
					Title:       "First Task",
					Description: "First description",
					DueDate:     parseTime("2024-09-05T15:04:05+05:00"),
					CreatedAt:   parseTime("2022-09-05T15:04:05+05:00"),
					UpdatedAt:   parseTime("2023-09-05T15:04:05+05:00"),
				}, nil)
			},
			expectedCode: 200,
			expectedBody: `{
								"id": 129,
								"title": "First Task",
								"description": "First description",
								"due_date": "2024-09-05T15:04:05+05:00",
								"created_at": "2022-09-05T15:04:05+05:00",
								"updated_at": "2023-09-05T15:04:05+05:00"
							}`,
		},
		{
			name:       "400_invalid_content_type",
			inputParam: "129",
			inputId:    129,
			inputBody: `{
							"title": "First Task",
							"description": "First description",
							"due_date": "2024-09-05T15:04:05+05:00"
						}`,
			inputDTO:      &dto.TaskUpdate{},
			inputContType: "text/plain",
			mockBehaviour: func(s *mockservice.MockITaskService, id int, update *dto.TaskUpdate) {},
			expectedCode:  400,
			expectedBody:  `{"error":"content-type is not application/json"}`,
		},
		{
			name:       "400_invalid_param",
			inputParam: "129f",
			inputId:    129,
			inputBody: `{
							"title": "First Task",
							"description": "First description",
							"due_date": "2024-09-05T15:04:05+05:00"
						}`,
			inputDTO:      &dto.TaskUpdate{},
			inputContType: "application/json",
			mockBehaviour: func(s *mockservice.MockITaskService, id int, update *dto.TaskUpdate) {},
			expectedCode:  400,
			expectedBody:  `{"error":"strconv.Atoi: parsing \"129f\": invalid syntax"}`,
		},
		{
			name:       "500_unknown_error",
			inputParam: "129",
			inputId:    129,
			inputBody: `{
							"title": "First Task",
							"description": "First description",
							"due_date": "2024-09-05T15:04:05+05:00"
						}`,
			inputDTO: &dto.TaskUpdate{
				Title:       "First Task",
				Description: "First description",
				DueDate:     parseTime("2024-09-05T15:04:05+05:00"),
			},
			inputContType: "application/json",
			mockBehaviour: func(s *mockservice.MockITaskService, id int, update *dto.TaskUpdate) {
				s.EXPECT().UpdateById(id, update).Return(&dto.TaskRead{}, errors.New("some error"))
			},
			expectedCode: 500,
			expectedBody: `{"error":"some error"}`,
		},
		{
			name:       "400_invalid_body_struct",
			inputParam: "129",
			inputId:    129,
			inputBody: `{
							"title": "First Task",
							"description": "First description",
							"due_date": 2024-09-05T15:04:05+05:00"
						}`,
			inputDTO:      &dto.TaskUpdate{},
			inputContType: "application/json",
			mockBehaviour: func(s *mockservice.MockITaskService, id int, update *dto.TaskUpdate) {},
			expectedCode:  400,
			expectedBody:  `{"error":"invalid character '-' after object key:value pair"}`,
		},
		{
			name:       "400_invalid_all_values",
			inputParam: "129",
			inputId:    129,
			inputBody: `{
							"title": "",
							"description": "",
							"due_date": "2024-09-05T15:04:75+05:00"
						}`,
			inputDTO:      &dto.TaskUpdate{},
			inputContType: "application/json",
			mockBehaviour: func(s *mockservice.MockITaskService, id int, update *dto.TaskUpdate) {},
			expectedCode:  400,
			expectedBody:  `{"error":"Title is required;Description is required;DueDate is required and must be in RFC3339 format;"}`,
		},
		{
			name:       "400_invalid_date_value",
			inputParam: "129",
			inputId:    129,
			inputBody: `{
							"title": "First Task",
							"description": "First description",
							"due_date": "2024-09-05T15:04:75+05:00"
						}`,
			inputDTO:      &dto.TaskUpdate{},
			inputContType: "application/json",
			mockBehaviour: func(s *mockservice.MockITaskService, id int, update *dto.TaskUpdate) {},
			expectedCode:  400,
			expectedBody:  `{"error":"DueDate is required and must be in RFC3339 format;"}`,
		},
		{
			name:       "404_no_rows_found",
			inputParam: "129",
			inputId:    129,
			inputBody: `{
							"title": "First Task",
							"description": "First description",
							"due_date": "2024-09-05T15:04:05+05:00"
						}`,
			inputDTO: &dto.TaskUpdate{
				Title:       "First Task",
				Description: "First description",
				DueDate:     parseTime("2024-09-05T15:04:05+05:00"),
			},
			inputContType: "application/json",
			mockBehaviour: func(s *mockservice.MockITaskService, id int, update *dto.TaskUpdate) {
				s.EXPECT().UpdateById(id, update).Return(&dto.TaskRead{}, pgx.ErrNoRows)
			},
			expectedCode: 404,
			expectedBody: `{"error":"no rows in result set"}`,
		},
		{
			name:       "500_unknown_error",
			inputParam: "129",
			inputId:    129,
			inputBody: `{
							"title": "First Task",
							"description": "First description",
							"due_date": "2024-09-05T15:04:05+05:00"
						}`,
			inputDTO: &dto.TaskUpdate{
				Title:       "First Task",
				Description: "First description",
				DueDate:     parseTime("2024-09-05T15:04:05+05:00"),
			},
			inputContType: "application/json",
			mockBehaviour: func(s *mockservice.MockITaskService, id int, update *dto.TaskUpdate) {
				s.EXPECT().UpdateById(id, update).Return(&dto.TaskRead{}, errors.New("some error"))
			},
			expectedCode: 500,
			expectedBody: `{"error":"some error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()
			taskService := mockservice.NewMockITaskService(c)
			testCase.mockBehaviour(taskService, testCase.inputId, testCase.inputDTO)

			services := service.Services{Task: taskService}
			handler := NewHandler(Deps{
				Service: services,
				Logger:  logging.GetLoggerTest(),
			})

			//Test server
			r := httprouter.New()
			r.PUT("/tasks/:id", handler.taskUpdateById)

			//http test
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/tasks/"+testCase.inputParam, strings.NewReader(testCase.inputBody))
			req.Header.Set("Content-Type", testCase.inputContType)

			//Perform request
			r.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, testCase.expectedCode, w.Code)
			if testCase.expectedBody != "" {
				if testCase.expectedBody == "nil" {
					testCase.expectedBody = ""
				}
				assert.JSONEq(t, testCase.expectedBody, w.Body.String())
			}
			if testCase.bodyMustContain != "" {
				assert.Contains(t, w.Body.String(), testCase.bodyMustContain)
			}
		})
	}
}

func TestHandler_taskDeleteById(t *testing.T) {
	type mockBehaviour func(s *mockservice.MockITaskService, id int)

	testTable := []struct {
		name            string
		inputParam      string
		inputId         int
		mockBehaviour   mockBehaviour
		expectedCode    int
		expectedBody    string
		bodyMustContain string
	}{
		{
			name:       "204_valid_param",
			inputParam: "129",
			inputId:    129,
			mockBehaviour: func(s *mockservice.MockITaskService, id int) {
				s.EXPECT().DeleteById(id).Return(nil)
			},
			expectedCode: 204,
			expectedBody: "",
		},
		{
			name:       "404_no_rows_found",
			inputParam: "129",
			inputId:    129,
			mockBehaviour: func(s *mockservice.MockITaskService, id int) {
				s.EXPECT().DeleteById(id).Return(pgx.ErrNoRows)
			},
			expectedCode: 404,
			expectedBody: "",
		},
		{
			name:       "500_unknown_error",
			inputParam: "129",
			inputId:    129,
			mockBehaviour: func(s *mockservice.MockITaskService, id int) {
				s.EXPECT().DeleteById(id).Return(errors.New("some error"))
			},
			expectedCode: 404,
			expectedBody: `{"error":"some error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()
			taskService := mockservice.NewMockITaskService(c)
			testCase.mockBehaviour(taskService, testCase.inputId)

			services := service.Services{Task: taskService}
			handler := NewHandler(Deps{
				Service: services,
				Logger:  logging.GetLoggerTest(),
			})

			//Test server
			r := httprouter.New()
			r.DELETE("/tasks/:id", handler.taskDeleteById)

			//http test
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/tasks/"+testCase.inputParam, strings.NewReader(""))

			//Perform request
			r.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, testCase.expectedCode, w.Code)
			if testCase.expectedBody != "" {
				if testCase.expectedBody == "nil" {
					testCase.expectedBody = ""
				}
				assert.JSONEq(t, testCase.expectedBody, w.Body.String())
			}
			if testCase.bodyMustContain != "" {
				assert.Contains(t, w.Body.String(), testCase.bodyMustContain)
			}
		})
	}
}
