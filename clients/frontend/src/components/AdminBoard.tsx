import { useState } from "react";
import auction from "../assets/auction.png";
import { AssignedAuctions } from "./AssignedAuctions";
import { ClosedAuctions } from "./ClosedAuctions";
import { InProgressAuctions } from "./InProgressAuctions";
import { PendingAuctions } from "./PendingAuctions";

type Props = {};

export const AdminBoard: React.FC<Props> = () => {
  const [isPendingAuctionsOpen, setPendingAuctionsOpen] = useState(false);
  const [isAssignedAuctionsOpen, setAssignedOpen] = useState(false);
  const [isInProgressAuctionsOpen, setInProgressOpen] = useState(false);
  const [isClosedAuctionsOpen, setClosedOpen] = useState(false);

  const togglePendingAuctions = () => {
    setPendingAuctionsOpen(!isPendingAuctionsOpen);
  };

  const toggleAssigned = () => {
    setAssignedOpen(!isAssignedAuctionsOpen);
  };

  const toggleInProgress = () => {
    setInProgressOpen(!isInProgressAuctionsOpen);
  };

  const toggleClosed = () => {
    setClosedOpen(!isClosedAuctionsOpen);
  };

  return (
    <div className="container">
      <div onClick={togglePendingAuctions} style={{ cursor: "pointer" }}>
        <img
          style={{ width: "50px", height: "60px" }}
          src={auction}
          alt="profile-img"
        />
        <span>
          <strong>Pending Auctions</strong>
        </span>
      </div>
      <div style={{ display: isPendingAuctionsOpen ? "block" : "none" }}>
        Choose an auction to resolve the winning bid!
        <PendingAuctions />
      </div>
      <br />
      <div onClick={toggleAssigned} style={{ cursor: "pointer" }}>
        <img
          style={{ width: "50px", height: "60px" }}
          src={auction}
          alt="profile-img"
        />
        <span>
          <strong>Assigned Auctions</strong>
        </span>
      </div>
      <div style={{ display: isAssignedAuctionsOpen ? "block" : "none" }}>
        Choose an auction to update the status!
        <AssignedAuctions />
      </div>
      <br />
      <div onClick={toggleInProgress} style={{ cursor: "pointer" }}>
        <img
          style={{ width: "50px", height: "60px" }}
          src={auction}
          alt="profile-img"
        />
        <span>
          <strong>In Progress Auctions</strong>
        </span>
      </div>
      <div style={{ display: isInProgressAuctionsOpen ? "block" : "none" }}>
        Choose an auction to close it!
        <InProgressAuctions />
      </div>
      <br />
      <div onClick={toggleClosed} style={{ cursor: "pointer" }}>
        <img
          style={{ width: "50px", height: "60px" }}
          src={auction}
          alt="profile-img"
        />
        <span>
          <strong>Closed Auctions</strong>
        </span>
      </div>
      <div style={{ display: isClosedAuctionsOpen ? "block" : "none" }}>
        All the closed auction from residents!
        <ClosedAuctions />
      </div>
    </div>
  );
};
