import React, { useState, useEffect } from "react";
import { Loader } from "../common/loader/Loader";
import { getLoggedInUserDetails } from "../services/auth";

export const Home: React.FC = () => {
  const [content, setContent] = useState<string>("");

  useEffect(() => {
    getLoggedInUserDetails().then(
      (response) => {
        const userStr = JSON.stringify(response.data);
        localStorage.setItem("user", userStr);
        setContent(userStr);
      },
      (error) => {
        setContent("");
      }
    );
  }, []);

  return (
    <div className="container" style={{ textAlign: "center" }}>
      <header className="jumbotron">
        <h3>{content === "" ? <Loader /> : content}</h3>
      </header>
    </div>
  );
};
