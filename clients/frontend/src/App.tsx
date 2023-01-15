import React from "react";
import { useState, useEffect } from "react";
import { Routes, Route, Link, NavigateFunction, useNavigate } from "react-router-dom";
import "bootstrap/dist/css/bootstrap.min.css";
import "./App.css";

import * as AuthService from "./services/auth";
import { User } from './types/user';

import { Login } from "./components/Login";
// import Register from "./components/Register";
import { Home } from "./components/Home";
import { Requests } from "./components/Requests";
// import Profile from "./components/Profile";
// import BoardUser from "./components/BoardUser";
// import BoardModerator from "./components/BoardModerator";
// import BoardAdmin from "./components/BoardAdmin";

// import { EventBus } from "./common/EventBus";

const App: React.FC = () => {
  let navigate: NavigateFunction = useNavigate();
  // const [showModeratorBoard, setShowModeratorBoard] = useState<boolean>(false);
  // const [showAdminBoard, setShowAdminBoard] = useState<boolean>(false);
  const [currentUser, setCurrentUser] = useState<User | undefined>(undefined);

  useEffect(() => {
    const user = AuthService.getCurrentUser();

    if (user) {
      setCurrentUser(user);
      // setShowModeratorBoard(user.roles.includes("ROLE_MODERATOR"));
      // setShowAdminBoard(user.roles.includes("ROLE_ADMIN"));
    }
  }, []);

  const logOut = async () => {
    await AuthService.logout();
    setCurrentUser(undefined);
    navigate("/home");
    window.location.reload();
  };

  return (
    <div>
      <nav className="navbar navbar-expand navbar-dark bg-dark">
        <Link to={"/"} className="navbar-brand">
          Service Bid Platform
        </Link>
        <div className="navbar-nav mr-auto">
          <li className="nav-item">
            <Link to={"/home"} className="nav-link">
              Home
            </Link>
          </li>

          {/* {showModeratorBoard && (
            <li className="nav-item">
              <Link to={"/mod"} className="nav-link">
                Moderator Board
              </Link>
            </li>
          )}

          {showAdminBoard && (
            <li className="nav-item">
              <Link to={"/admin"} className="nav-link">
                Admin Board
              </Link>
            </li>
          )} */}

        </div>

        {currentUser ? (
          <div className="navbar-nav ml-auto">
            <li className="nav-item">
              <a href="/requests" className="nav-link">
                MyRequests
              </a>
            </li>
            <li className="nav-item">
              <a href="/logout" className="nav-link" onClick={logOut}>
                Logout
              </a>
            </li>
          </div>
        ) : (
          <div className="navbar-nav ml-auto">
            <li className="nav-item">
              <Link to={"/login"} className="nav-link">
                Login
              </Link>
            </li>
          </div>
        )}
      </nav>

      <div className="container mt-3">
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/home" element={<Home />} />
          <Route path="/login" element={<Login />} />
          <Route path="/requests" element={<Requests />} />
          {/*<Route path="/register" element={<Register />} />
          <Route path="/profile" element={<Profile />} />
          <Route path="/user" element={<BoardUser />} />
          <Route path="/mod" element={<BoardModerator />} />
          <Route path="/admin" element={<BoardAdmin />} /> */}
        </Routes>
      </div>
    </div>
  );
};

export default App;