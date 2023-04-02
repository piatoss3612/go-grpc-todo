import Row from "react-bootstrap/esm/Row";
import { Message } from "./Message";
import { useContext } from "react";
import { TodoContext } from "../context/context";

export const Messages = () => {
  const { requestType, request, response } = useContext(TodoContext);

  return (
    <Row className="mt-3">
      <Message
        header="Request sent"
        title={`${requestType.toUpperCase()} Request`}
        content={request}
      />
      <Message
        header="Response received"
        title={`${requestType.toUpperCase()} Response`}
        content={response}
      />
    </Row>
  );
};
