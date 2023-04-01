import { useState } from "react";
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";

const requestTypes = ["gRPC", "http"];
const endpoints = ["v1/todo", "v1/todo/all"];
const Methods = ["GET", "POST", "PUT", "DELETE"];
const priorities = [
  { value: 1, label: "LOW" },
  { value: 2, label: "MEDIUM" },
  { value: 3, label: "HIGH" },
];

export const RequestForm = () => {
  const [requestType, setRequestType] = useState(requestTypes[0]);
  const [streamRepeat, setStreamRepeat] = useState(0);
  const [method, setMethod] = useState(Methods[0]);
  const [endpoint, setEndpoint] = useState(endpoints[0]);
  const [todoId, setTodoId] = useState("");
  const [todoContent, setTodoContent] = useState("");
  const [todoPriority, setTodoPriority] = useState(priorities[0].value);
  const [todoCompleted, setTodoCompleted] = useState(false);

  const handleRequestTypeChange = (e) => {
    setRequestType(e.target.value);
  };

  const handleStreamRepeatChange = (e) => {
    setStreamRepeat(e.target.value);
  };

  const handleMethodChange = (e) => {
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
    console.log("submit");
    setRequestType(requestTypes[0]);
    setStreamRepeat(0);
    setMethod(Methods[0]);
    setEndpoint(endpoints[0]);
    setTodoId("");
    setTodoContent("");
    setTodoPriority(priorities[0].value);
    setTodoCompleted(false);
  };

  return (
    <Form onSubmit={handleSubmit}>
      <Form.Group controlId="formRequestType" className="mb-2">
        <Form.Label>Request Type</Form.Label>
        <Form.Control
          as="select"
          value={requestType}
          onChange={handleRequestTypeChange}
        >
          {requestTypes.map((requestType) => (
            <option key={requestType}>{requestType}</option>
          ))}
        </Form.Control>
      </Form.Group>
      {requestType === "gRPC" && (
        <Form.Group controlId="formStreamRepeat" className="mb-2">
          <Form.Label>Stream Repeat</Form.Label>
          <Form.Range
            step={1}
            min={0}
            max={10}
            value={streamRepeat}
            onChange={handleStreamRepeatChange}
          />
        </Form.Group>
      )}
      <Form.Group controlId="formMethod" className="mb-2">
        <Form.Label>Method</Form.Label>
        <Form.Control as="select" value={method} onChange={handleMethodChange}>
          {Methods.map((method) => (
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
            {endpoints.map((endpoint) => (
              <option key={endpoint}>{endpoint}</option>
            ))}
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
