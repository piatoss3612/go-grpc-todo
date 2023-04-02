import { Container } from "react-bootstrap";
import { Messages } from "./components/Messages";
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";
import { RequestForm } from "./components/RequestForm";

const App = () => {
  return (
    <Container>
      <Row>
        <Col>
          <h1 className="mt-5">Todo gRPC Test</h1>
          <hr />
          <RequestForm />
        </Col>
      </Row>
      <Messages />
    </Container>
  );
};

export default App;
