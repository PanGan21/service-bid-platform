import { useState } from "react";
import request from "../assets/request.png";
import { PendingRequests } from "./PendingRequests";

type Props = {};

export const AdminBoard: React.FC<Props> = () => {
  const [isPendingRequestsOpen, setPendingRequestsOpen] = useState(false);
  const [isInProgressRequestsOpen, setInProgressRequestsOpen] = useState(false);

  const togglePendingRequests = () => {
    setPendingRequestsOpen(!isPendingRequestsOpen);
  };

  const toggleInProgressRequests = () => {
    setInProgressRequestsOpen(!isInProgressRequestsOpen);
  };

  return (
    <div className="container">
      <div onClick={togglePendingRequests} style={{ cursor: "pointer" }}>
        <img
          style={{ width: "50px", height: "60px" }}
          src={request}
          alt="profile-img"
        />
        <span>
          <strong>Pending Requests</strong>
        </span>
      </div>
      <div style={{ display: isPendingRequestsOpen ? "block" : "none" }}>
        Choose a request to resolve the winning bid!
        <PendingRequests />
      </div>
      <br />
      <div onClick={toggleInProgressRequests} style={{ cursor: "pointer" }}>
        <img
          style={{ width: "50px", height: "60px" }}
          src={request}
          alt="profile-img"
        />
        <span>
          <strong>In Progress Requests</strong>
        </span>
      </div>
      <div style={{ display: isInProgressRequestsOpen ? "block" : "none" }}>
      </div>
    </div>
  );
};
