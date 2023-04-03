import { createContext, useState } from "react";

export const TodoContext = createContext();

export const TodoProvider = ({ children }) => {
  const requestType = "http";
  const endpoints = ["v1/todo", "v1/todo/all"];
  const methods = ["GET", "POST", "PUT", "DELETE"];
  const priorities = [
    { value: 1, label: "LOW" },
    { value: 2, label: "MEDIUM" },
    { value: 3, label: "HIGH" },
  ];

  const [request, setRequest] = useState("Nothing sent yet");
  const [response, setResponse] = useState("Nothing received yet");

  const todoRequestBaseUrl = "http://localhost:80/api";

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
      const response = await fetch(
        `${todoRequestBaseUrl}/${endpoint}`,
        request
      );
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
        ? `${todoRequestBaseUrl}/${endpoint}/${id}`
        : `${todoRequestBaseUrl}/${endpoint}`;
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
      ? `${todoRequestBaseUrl}/${endpoint}/${id}`
      : `${todoRequestBaseUrl}/${endpoint}`;

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
        ? `${todoRequestBaseUrl}/${endpoint}/${id}`
        : `${todoRequestBaseUrl}/${endpoint}`;
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
        methods,
        endpoints,
        priorities,
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
