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
import { Home } from "./components/Home";
import { OpenAuctions } from "./components/OpenAuctions";
import { CreateServiceRequest } from "./components/CreateServiceRequest";
import { CreateBid } from "./components/CreateBid";
import { Footer } from "./components/Footer";
import { AdminBoard } from "./components/AdminBoard";
import { AssignedAuction } from "./components/AssignedAuction";
import { Assignments } from "./components/Assignments";
import { NewServiceRequests } from "./components/NewServiceRequests";
import { UpdateRequestStatus } from "./components/UpdateRequestStatus";
import { UpdateAuctionStatus } from "./components/UpdateAuctionStatus";

const App: React.FC = () => {
  let navigate: NavigateFunction = useNavigate();
  const [showAdminBoard, setShowAdminBoard] = useState<boolean>(false);
  const [showResidentBoard, setShowResidentBoard] = useState<boolean>(false);
  const [showBidderBoard, setShowBidderBoard] = useState<boolean>(false);
  const [currentUser, setCurrentUser] = useState<User | undefined>(undefined);

  useEffect(() => {
    const user = AuthService.getCurrentUser();

    if (user) {
      setCurrentUser(user);
      setShowAdminBoard(user.Roles.includes("ADMIN"));
      setShowBidderBoard(user.Roles.includes("BIDDER"));
      setShowResidentBoard(user.Roles.includes("RESIDENT"));
    }
    if (!user) {
      navigate("/login");
    }
  }, [navigate]);

  const logOut = async () => {
    await AuthService.logout();
    setCurrentUser(undefined);
    navigate("/home");
    window.location.reload();
  };

  const allowRoles = (predicates: boolean[]): boolean => {
    const isAtLeastOneRolePresent = predicates.some((item) => item);
    if (isAtLeastOneRolePresent) return true;
    return false;
  };

  return (
    <div>
      <nav className="navbar navbar-expand navbar-dark bg-dark">
        <Link to={"/home"} className="navbar-brand">
          Service Bid Platform
        </Link>
        <div className="navbar-nav mr-auto">
          <li className="nav-item">
            <Link to={"/home"} className="nav-link">
              Home
            </Link>
          </li>
          {allowRoles([showAdminBoard]) && (
            <li className="nav-item">
              <Link to={"/admin"} className="nav-link">
                Admin Board
              </Link>
            </li>
          )}
          {allowRoles([showAdminBoard]) && (
            <li className="nav-item">
              <Link to={"/new-service-requests"} className="nav-link">
                New Service Requests
              </Link>
            </li>
          )}
        </div>

        {currentUser ? (
          <div className="navbar-nav ml-auto">
            {allowRoles([showBidderBoard, showAdminBoard]) && (
              <li className="nav-item">
                <a href="/open-auctions" className="nav-link">
                  Open Auctions
                </a>
              </li>
            )}
            {allowRoles([showResidentBoard, showAdminBoard]) && (
              <li className="nav-item">
                <a href="/create-service-request" className="nav-link">
                  Create Service Request
                </a>
              </li>
            )}
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
          <Route path="/open-auctions" element={<OpenAuctions />} />
          <Route
            path="/new-service-requests"
            element={<NewServiceRequests />}
          />
          <Route
            path="/create-service-request"
            element={<CreateServiceRequest />}
          />
          <Route path="/new-bid" element={<CreateBid />} />
          <Route path="/admin" element={<AdminBoard />} />
          <Route path="/assigned-auction" element={<AssignedAuction />} />
          <Route
            path="/update-auction-status"
            element={<UpdateAuctionStatus />}
          />
          <Route
            path="/update-request-status"
            element={<UpdateRequestStatus />}
          />
          <Route path="/assignments" element={<Assignments />} />
        </Routes>
      </div>
      <Footer />
    </div>
  );
};

export default App;
