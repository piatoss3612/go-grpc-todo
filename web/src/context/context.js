import { createContext, useState } from "react";

export const TodoContext = createContext();

export const TodoProvider = ({ children }) => {
  const [requestType, setRequestType] = useState("");
  const [request, setRequest] = useState("Nothing sent yet");
  const [response, setResponse] = useState("Nothing received yet");

  const httpHost = "http://localhost:8080";

  const todoHttpAddRequest = async (endpoint, content, priority) => {
    const request = {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        content: content,
        priority: parseInt(priority, 10),
      }),
    };

    try {
      const response = await fetch(`${httpHost}/${endpoint}`, request);
      const data = await response.json();
      setRequest(JSON.stringify(request));
      setResponse(JSON.stringify(data));
    } catch (error) {
      setRequest(JSON.stringify(request));
      setResponse(error.message);
    }
  };

  const todoHttpGetRequest = async (endpoint, id) => {
    const request = id ? `id: ${id}` : "all";
    try {
      const url = id
        ? `${httpHost}/${endpoint}/${id}`
        : `${httpHost}/${endpoint}`;
      const response = await fetch(url);
      const data = await response.text();
      setRequest(request);
      setResponse(data);
    } catch (error) {
      setRequest(request);
      setResponse(error.message);
    }
  };

  const todoHttpUpdateRequest = async (
    endpoint,
    id,
    content,
    priority,
    isDone
  ) => {
    const request = {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        id: id,
        content: content,
        priority: parseInt(priority, 10),
        isDone: isDone,
      }),
    };

    const url = id
      ? `${httpHost}/${endpoint}/${id}`
      : `${httpHost}/${endpoint}`;

    try {
      const response = await fetch(url, request);
      const data = await response.json();
      setRequest(JSON.stringify(request));
      setResponse(JSON.stringify(data));
    } catch (error) {
      setRequest(JSON.stringify(request));
      setResponse(error.message);
    }
  };

  const todoHttpDeleteRequest = async (endpoint, id) => {
    const request = id ? `id: ${id}` : "all";
    try {
      const url = id
        ? `${httpHost}/${endpoint}/${id}`
        : `${httpHost}/${endpoint}`;
      const response = await fetch(url, { method: "DELETE" });
      const data = await response.json();
      setRequest(request);
      setResponse(JSON.stringify(data));
    } catch (error) {
      setRequest(request);
      setResponse(error.message);
    }
  };

  return (
    <TodoContext.Provider
      value={{
        requestType,
        setRequestType,
        request,
        setRequest,
        response,
        setResponse,
        todoHttpGetRequest,
        todoHttpAddRequest,
        todoHttpUpdateRequest,
        todoHttpDeleteRequest,
      }}
    >
      {children}
    </TodoContext.Provider>
  );
};
