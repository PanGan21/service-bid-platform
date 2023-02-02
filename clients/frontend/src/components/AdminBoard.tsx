import { useState } from "react";
import request from "../assets/request.png";
import { AssignedRequests } from "./AssignedRequests";
import { ClosedRequests } from "./ClosedRequests";
import { InProgressRequests } from "./InProgressRequests";
import { PendingRequests } from "./PendingRequests";

type Props = {};

export const AdminBoard: React.FC<Props> = () => {
  const [isPendingRequestsOpen, setPendingRequestsOpen] = useState(false);
  const [isAssignedRequestsOpen, setAssignedOpen] = useState(false);
  const [isInProgressRequestsOpen, setInProgressOpen] = useState(false);
  const [isClosedRequestsOpen, setClosedOpen] = useState(false);

  const togglePendingRequests = () => {
    setPendingRequestsOpen(!isPendingRequestsOpen);
  };

  const toggleAssigned = () => {
    setAssignedOpen(!isAssignedRequestsOpen);
  };

  const toggleInProgress = () => {
    setInProgressOpen(!isInProgressRequestsOpen);
  };

  const toggleClosed = () => {
    setClosedOpen(!isClosedRequestsOpen);
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
      <div onClick={toggleAssigned} style={{ cursor: "pointer" }}>
        <img
          style={{ width: "50px", height: "60px" }}
          src={request}
          alt="profile-img"
        />
        <span>
          <strong>Assigned Requests</strong>
        </span>
      </div>
      <div style={{ display: isAssignedRequestsOpen ? "block" : "none" }}>
        Choose a request to update the status!
        <AssignedRequests />
      </div>
      <br />
      <div onClick={toggleInProgress} style={{ cursor: "pointer" }}>
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
        Choose a request to close it!
        <InProgressRequests />
      </div>
      <br />
      <div onClick={toggleClosed} style={{ cursor: "pointer" }}>
        <img
          style={{ width: "50px", height: "60px" }}
          src={request}
          alt="profile-img"
        />
        <span>
          <strong>Closed Requests</strong>
        </span>
      </div>
      <div style={{ display: isClosedRequestsOpen ? "block" : "none" }}>
        All the closed request from residents!
        <ClosedRequests />
      </div>
    </div>
  );
};
