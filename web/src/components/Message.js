import Col from "react-bootstrap/esm/Col";
import Card from "react-bootstrap/esm/Card";

export const Message = ({ header, title, content }) => {
  return (
    <Col>
      <Card>
        <Card.Header>{header}</Card.Header>
        <Card.Body>
          <Card.Title className="mb-2">{title}</Card.Title>
          <Card.Text>{content}</Card.Text>
        </Card.Body>
      </Card>
    </Col>
  );
};
