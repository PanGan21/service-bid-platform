import React from "react";
import { useState, useEffect } from "react";
import {
  Routes,
  Route,
  Link,
  NavigateFunction,
  useNavigate,
} from "react-router-dom";
import "bootstrap/dist/css/bootstrap.min.css";
import "./App.css";

import * as AuthService from "./services/auth";
import { User } from "./types/user";

import { Login } from "./components/Login";
// import Register from "./components/Register";
import { Home } from "./components/Home";
import { AllRequests } from "./components/AllRequests";
import { CreateRequest } from "./components/CreateRequest";
import { CreateBid } from "./components/CreateBid";
import { Footer } from "./components/Footer";
import { AdminBoard } from "./components/AdminBoard";
import { AssignRequest } from "./components/AssignRequest";
import { UpdateRequestStatus } from "./components/UpdateRequestStatus";
import { Assignments } from "./components/Assignments";
// import Profile from "./components/Profile";
// import BoardUser from "./components/BoardUser";
// import BoardModerator from "./components/BoardModerator";
// import BoardAdmin from "./components/BoardAdmin";

const App: React.FC = () => {
  let navigate: NavigateFunction = useNavigate();
  const [showAdminBoard, setShowAdminBoard] = useState<boolean>(false);
  const [currentUser, setCurrentUser] = useState<User | undefined>(undefined);

  useEffect(() => {
    const user = AuthService.getCurrentUser();

    if (user) {
      setCurrentUser(user);
      setShowAdminBoard(user.Roles.includes("ADMIN"));
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
          {showAdminBoard && (
            <li className="nav-item">
              <Link to={"/admin"} className="nav-link">
                Admin Board
              </Link>
            </li>
          )}
        </div>

        {currentUser ? (
          <div className="navbar-nav ml-auto">
            <li className="nav-item">
              <a href="/all-requests" className="nav-link">
                AllRequests
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
          <Route path="/all-requests" element={<AllRequests />} />
          <Route path="/new-request" element={<CreateRequest />} />
          <Route path="/new-bid" element={<CreateBid />} />
          <Route path="/admin" element={<AdminBoard />} />
          <Route path="/assign-request" element={<AssignRequest />} />
          <Route
            path="/update-request-status"
            element={<UpdateRequestStatus />}
          />
          <Route path="/assignments" element={<Assignments />} />
          {/*<Route path="/register" element={<Register />} />
          <Route path="/profile" element={<Profile />} />
          <Route path="/user" element={<BoardUser />} />
          <Route path="/mod" element={<BoardModerator />} />
          <Route path="/admin" element={<BoardAdmin />} /> */}
        </Routes>
      </div>
      <Footer />
    </div>
  );
};

export default App;
