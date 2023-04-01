import { Container } from "react-bootstrap";
import { RequestGenerator } from "./components/RequestGenerator";
import { Messages } from "./components/Messages";

const App = () => {
  return (
    <Container>
      <RequestGenerator />
      <Messages />
    </Container>
  );
};

export default App;
