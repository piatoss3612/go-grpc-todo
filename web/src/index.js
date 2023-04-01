import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App";
import "./index.css";

import "bootstrap/dist/css/bootstrap.min.css";

import { TodoProvider } from "./context/context";

const root = ReactDOM.createRoot(document.getElementById("root"));
root.render(
  <TodoProvider>
    <App />
  </TodoProvider>
);
