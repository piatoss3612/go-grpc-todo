import Col from "react-bootstrap/esm/Col";
import Card from "react-bootstrap/esm/Card";
import { useContext } from "react";
import { TodoContext } from "../context/context";

export const ResponseMessage = () => {
  const { requestType, response } = useContext(TodoContext);

  return (
    <Col>
      <Card>
        <Card.Header>Response received</Card.Header>
        <Card.Body>
          <Card.Title className="mb-2">
            {requestType ? `${requestType} Response` : "Response"}
          </Card.Title>
          <Card.Text>{response}</Card.Text>
        </Card.Body>
      </Card>
    </Col>
  );
};
