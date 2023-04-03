import { useContext, useState } from "react";
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";
import { TodoContext } from "../context/context";

export const RequestForm = () => {
  const {
    requestType,
    methods,
    endpoints,
    priorities,
    todoHttpGetRequest,
    todoHttpAddRequest,
    todoHttpUpdateRequest,
    todoHttpDeleteRequest,
  } = useContext(TodoContext);

  const [method, setMethod] = useState(methods[0]);
  const [endpoint, setEndpoint] = useState(endpoints[0]);
  const [todoId, setTodoId] = useState("");
  const [todoContent, setTodoContent] = useState("");
  const [todoPriority, setTodoPriority] = useState(priorities[0].value);
  const [todoCompleted, setTodoCompleted] = useState(false);

  const handleMethodChange = (e) => {
    if (e.target.value === "POST" || e.target.value === "PUT") {
      setEndpoint(endpoints[0]);
    }
    setMethod(e.target.value);
  };

  const handleEndpointChange = (e) => {
    setEndpoint(e.target.value);
  };

  const handleTodoIdChange = (e) => {
    setTodoId(e.target.value);
  };

  const handleTodoContentChange = (e) => {
    setTodoContent(e.target.value);
  };

  const handleTodoPriorityChange = (e) => {
    setTodoPriority(e.target.value);
  };

  const toggleCompleted = () => {
    setTodoCompleted((prev) => !prev);
  };

  const handleSubmit = (e) => {
    e.preventDefault();

    try {
      switch (requestType) {
        case "http":
          switch (method) {
            case "GET":
              todoHttpGetRequest(endpoint, todoId);
              break;
            case "POST":
              todoHttpAddRequest(endpoint, todoContent, todoPriority);
              break;
            case "PUT":
              todoHttpUpdateRequest(
                endpoint,
                todoId,
                todoContent,
                todoPriority,
                todoCompleted
              );
              break;
            case "DELETE":
              todoHttpDeleteRequest(endpoint, todoId);
              break;
            default:
              throw new Error("Unsupported method");
          }
          break;
        default:
          throw new Error("Unsupported request type");
      }
    } catch (error) {
      console.log(error);
      return;
    }

    setTodoId("");
    setTodoContent("");
    setTodoPriority(priorities[0].value);
  };

  return (
    <Form onSubmit={handleSubmit}>
      <Form.Group controlId="formRequestType" className="mb-2">
        <Form.Label>Request Type</Form.Label>
        <Form.Control type={"text"} value={requestType} disabled></Form.Control>
      </Form.Group>
      <Form.Group controlId="formMethod" className="mb-2">
        <Form.Label>Method</Form.Label>
        <Form.Control as="select" value={method} onChange={handleMethodChange}>
          {methods.map((method) => (
            <option key={method}>{method}</option>
          ))}
        </Form.Control>
      </Form.Group>
      {requestType === "http" && (
        <Form.Group controlId="formEndpoint" className="mb-2">
          <Form.Label>Endpoint</Form.Label>
          <Form.Control
            as="select"
            value={endpoint}
            onChange={handleEndpointChange}
          >
            <option>{endpoints[0]}</option>
            <option disabled={method === "POST" || method === "PUT"}>
              {endpoints[1]}
            </option>
          </Form.Control>
        </Form.Group>
      )}
      {method !== "POST" && (
        <Form.Group controlId="formTodoId" className="mb-2">
          <Form.Label>ID</Form.Label>
          <Form.Control
            value={todoId}
            onChange={handleTodoIdChange}
          ></Form.Control>
        </Form.Group>
      )}
      {(method === "POST" || method === "PUT") && (
        <Form.Group controlId="formTodoContent" className="mb-2">
          <Form.Label>Content</Form.Label>
          <Form.Control
            value={todoContent}
            onChange={handleTodoContentChange}
          ></Form.Control>
        </Form.Group>
      )}
      {(method === "POST" || method === "PUT") && (
        <Form.Group controlId="formTodoPriority" className="mb-2">
          <Form.Label>Priority</Form.Label>
          <Form.Control
            as="select"
            value={todoPriority}
            onChange={handleTodoPriorityChange}
          >
            {priorities.map((priority) => (
              <option key={priority.value} value={priority.value}>
                {priority.label}
              </option>
            ))}
          </Form.Control>
        </Form.Group>
      )}
      {method === "PUT" && (
        <Form.Group controlId="formTodoCompleted" className="mb-2">
          <Form.Check
            type="checkbox"
            value={todoCompleted}
            onChange={toggleCompleted}
            label="Completed"
          ></Form.Check>
        </Form.Group>
      )}

      <Button type="submit" className="mt-2">
        Submit
      </Button>
    </Form>
  );
};
