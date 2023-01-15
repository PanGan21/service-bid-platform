import React, { useState, useEffect } from "react";
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
                const _content = (error.response.data.error || error.message || JSON.stringify(error))
                setContent(_content);
            }
        );
    }, []);

    return (
        <div className="container">
            <header className="jumbotron">
                <h3>{content}</h3>
            </header>
        </div>
    );
};