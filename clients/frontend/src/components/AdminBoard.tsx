import { useState } from "react";
import request from "../assets/request.png";
import { PendingRequests } from "./PendingRequests";

type Props = {};

export const AdminBoard: React.FC<Props> = () => {
  const [isPendingRequestsOpen, setPendingRequestsOpen] = useState(false);

  const togglePendingRequests = () => {
    setPendingRequestsOpen(!isPendingRequestsOpen);
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
    </div>
  );
};
