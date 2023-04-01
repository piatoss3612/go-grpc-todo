import { createContext, useState } from "react";

export const TodoContext = createContext();

export const TodoProvider = ({ children }) => {
  const [todos, setTodos] = useState([]);
  const [requestType, setRequestType] = useState("");
  const [request, setRequest] = useState("Nothing sent yet");
  const [response, setResponse] = useState("Nothing received yet");

  return (
    <TodoContext.Provider
      value={{
        todos,
        setTodos,
        requestType,
        setRequestType,
        request,
        setRequest,
        response,
        setResponse,
      }}
    >
      {children}
    </TodoContext.Provider>
  );
};
