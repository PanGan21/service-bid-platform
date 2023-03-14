import React, { useState, useEffect } from "react";
import { Loader } from "../common/loader/Loader";
import { User } from "../types/user";
import { MyAuctions } from "./MyAuctions";
import auction from "../assets/auction.png";
import bid from "../assets/bid.png";
import { MyBids } from "./MyBids";
import { ProfileImageBadge } from "./ProfileImageBadge";
import { countOwnAssignments } from "../services/auction";
import { useNavigate } from "react-router-dom";
import { MyRejectedAuctions } from "./MyRejectedAuctions";

export const Home: React.FC = () => {
  const [user, setUser] = useState<User | undefined>(undefined);
  const [isMyAuctionsOpen, setMyAuctionsOpen] = useState(false);
  const [isMyBidsOpen, setMyBidsOpen] = useState(false);
  const [isRejectedAuctionsOpen, setIsRejectedAuctionsOpen] = useState(false);
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

  const toggleMyAuctions = () => {
    setMyAuctionsOpen(!isMyAuctionsOpen);
  };

  const toggleMyBids = () => {
    setMyBidsOpen(!isMyBidsOpen);
  };

  const toggleRejectedAuctions = () => {
    setIsRejectedAuctionsOpen(!isRejectedAuctionsOpen);
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
        <div>
          <h6>
            Email: <strong>{user.Email}</strong>
          </h6>
          <h6>
            Phone: <strong>{user.Phone}</strong>
          </h6>

          {user.Roles.length > 0 ? (
            <div>
              <strong>Roles:</strong>
              <ul>
                {user.Roles &&
                  user.Roles.map((role, index) => <li key={index}>{role}</li>)}
              </ul>
            </div>
          ) : (
            <></>
          )}
        </div>
      </div>
      <br />
      <div onClick={toggleMyAuctions} style={{ cursor: "pointer" }}>
        <img
          style={{ width: "50px", height: "60px" }}
          src={auction}
          alt="auctions-img"
        />
        <span>
          <strong>MyAuctions</strong>
        </span>
      </div>
      <div style={{ display: isMyAuctionsOpen ? "block" : "none" }}>
        <MyAuctions />
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
      <br />
      <div onClick={toggleRejectedAuctions} style={{ cursor: "pointer" }}>
        <img
          style={{ width: "50px", height: "60px" }}
          src={auction}
          alt="profile-img"
        />
        <span>
          <strong>My Rejected Auctions</strong>
        </span>
      </div>
      <div style={{ display: isRejectedAuctionsOpen ? "block" : "none" }}>
        <MyRejectedAuctions />
      </div>
    </div>
  );
};
