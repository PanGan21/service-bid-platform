import React, { useState, useEffect } from "react";
import { Loader } from "../common/loader/Loader";
import { User } from "../types/user";
import { MyRequests } from "./MyRequests";
import request from "../assets/request.png";
import bid from "../assets/bid.png";
import { MyBids } from "./MyBids";
import { ProfileImageBadge } from "./ProfileImageBadge";
import { countOwnAssignments } from "../services/request";
import { useNavigate } from "react-router-dom";

export const Home: React.FC = () => {
  const [user, setUser] = useState<User | undefined>(undefined);
  const [isMyRequestsOpen, setMyRequestsOpen] = useState(false);
  const [isMyBidsOpen, setMyBidsOpen] = useState(false);
  const [badgeNumber, setBadgeNumber] = useState(0);
  const navigate = useNavigate();

  useEffect(() => {
    const userStr = localStorage.getItem("user");
    if (userStr) {
      const localUser = JSON.parse(userStr);
      setUser(localUser);
    }
    countOwnAssignments().then((response) => {
      if (response && response.data) {
        setBadgeNumber(response.data);
      }
    });
  }, []);

  if (!user) {
    return (
      <div className="container" style={{ textAlign: "center" }}>
        <header className="jumbotron">
          <h3>
            <Loader />
          </h3>
        </header>
      </div>
    );
  }

  const toggleMyRequests = () => {
    setMyRequestsOpen(!isMyRequestsOpen);
  };

  const toggleMyBids = () => {
    setMyBidsOpen(!isMyBidsOpen);
  };

  const handleBadgeClick = () => {
    navigate("/assignments", { state: user });
  };

  return (
    <div className="container">
      <div className="card card-container" style={{ alignItems: "center" }}>
        <header className="jumbotron">
          <h3>
            <ProfileImageBadge
              src="//ssl.gstatic.com/accounts/ui/avatar_2x.png"
              badgeNumber={badgeNumber}
              onBadgeClick={handleBadgeClick}
            />
            MyProfile
          </h3>
        </header>
        <h4>
          Username: <strong>{user.Username}</strong>
        </h4>
        {user.Roles.length > 0 ? (
          <div>
            <strong>Authorities:</strong>
            <ul>
              {user.Roles &&
                user.Roles.map((role, index) => <li key={index}>{role}</li>)}
            </ul>
          </div>
        ) : (
          <></>
        )}
      </div>
      <br />
      <div onClick={toggleMyRequests} style={{ cursor: "pointer" }}>
        <img
          style={{ width: "50px", height: "60px" }}
          src={request}
          alt="requests-img"
        />
        <span>
          <strong>MyRequests</strong>
        </span>
      </div>
      <div style={{ display: isMyRequestsOpen ? "block" : "none" }}>
        <MyRequests />
      </div>
      <br />
      <div onClick={toggleMyBids} style={{ cursor: "pointer" }}>
        <img
          style={{ width: "50px", height: "60px" }}
          src={bid}
          alt="profile-img"
        />
        <span>
          <strong>MyBids</strong>
        </span>
      </div>
      <div style={{ display: isMyBidsOpen ? "block" : "none" }}>
        <MyBids />
      </div>
    </div>
  );
};
