import Col from "react-bootstrap/esm/Col";
import Card from "react-bootstrap/esm/Card";
import { useContext } from "react";
import { TodoContext } from "../context/context";

export const RequestMessage = () => {
  const { requestType, request } = useContext(TodoContext);

  return (
    <Col>
      <Card>
        <Card.Header>Request sent</Card.Header>
        <Card.Body>
          <Card.Title className="mb-2">
            {requestType ? `${requestType} Request` : "Request"}
          </Card.Title>
          <Card.Text>{request}</Card.Text>
        </Card.Body>
      </Card>
    </Col>
  );
};
