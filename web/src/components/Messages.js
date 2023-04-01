import Row from "react-bootstrap/esm/Row";
import { RequestMessage } from "./RequestMessage";
import { ResponseMessage } from "./ResponseMessage";

export const Messages = () => {
  return (
    <Row>
      <RequestMessage />
      <ResponseMessage />
    </Row>
  );
};
