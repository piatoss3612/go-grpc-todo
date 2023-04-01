import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";
import { RequestForm } from "./RequestForm";

export const RequestGenerator = () => {
  return (
    <Row>
      <Col>
        <h1 className="mt-5">Todo gRPC Test</h1>
        <hr />
        <RequestForm />
      </Col>
    </Row>
  );
};
