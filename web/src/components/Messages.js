import Row from "react-bootstrap/esm/Row";
import { RequestMessage } from "./RequestMessage";
import { ResponseMessage } from "./ResponseMessage";

export const Messages = () => {
  return (
    <Row className="mt-3">
      <RequestMessage />
      <ResponseMessage />
    </Row>
  );
};
