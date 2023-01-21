import React, { useState, useEffect } from "react";
import { Loader } from "../common/loader/Loader";
import { User } from "../types/user";
import { MyRequests } from "./MyRequests";
import request from "../assets/request.png";
import bid from "../assets/bid.png";
import { MyBids } from "./MyBids";

export const Home: React.FC = () => {
  const [user, setUser] = useState<User | undefined>(undefined);

  useEffect(() => {
    const userStr = localStorage.getItem("user");
    if (userStr) {
      const localUser = JSON.parse(userStr);
      setUser(localUser);
    }
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

  return (
    <div className="container">
      <header className="jumbotron">
        <h3>
          MyProfile
          <br />
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
      <br />
      <div>
        <img
          style={{ width: "50px", height: "60px" }}
          src={request}
          alt="profile-img"
        />
        <span>
          <strong>MyRequests</strong>
        </span>
      </div>
      <MyRequests />
      <br />
      <div>
        <img
          style={{ width: "50px", height: "60px" }}
          src={bid}
          alt="profile-img"
        />
        <span>
          <strong>MyBids</strong>
        </span>
      </div>
      <MyBids />
    </div>
  );
};
