import React from "react";
import ReactDOM from "react-dom";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import "./index.css";
import Login from "./Login";
import Register from "./Register";
import Temperatures from "./Temperatures";
import reportWebVitals from "./reportWebVitals";

ReactDOM.render(
  <Router>
    <Routes>
      <Route path="/temperatures" element={<Temperatures />} />
      <Route path="/register" element={<Register />} />
      <Route path="/" element={<Login />} />
    </Routes>
  </Router>,
  document.getElementById("root")
);

reportWebVitals();
