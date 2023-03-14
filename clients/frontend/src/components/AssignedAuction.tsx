import { useEffect, useState } from "react";
import { useLocation } from "react-router-dom";
import bid from "../assets/bid.png";
import { getUserDetailsById } from "../services/auth";
import { Auction } from "../types/auction";
import { UserDetails } from "../types/user";

type Props = {};

export const AssignedAuction: React.FC<Props> = () => {
  const { state }: { state: Auction } = useLocation();
  const [winnerDetails, setWinnerDetails] = useState<UserDetails | undefined>(
    undefined
  );

  useEffect(() => {
    if (state.WinnerId !== "") {
      getUserDetailsById(state.WinnerId).then((response) => {
        if (response.data) {
          setWinnerDetails(response.data);
        }
      });
    }
  }, [state.WinnerId]);

  return (
    <div className="col-md-12">
      <div className="card card-container">
        <img src={bid} alt="profile-img" className="profile-img-card" />
        <h3 style={{ textAlign: "center" }}>Winning bid</h3>
        <h6>
          <b>Auction Id:</b> {state.Id}
        </h6>
        <h6>
          <b>Bidder Id:</b> {state.WinnerId || "Pending"}
        </h6>
        <h6>
          <b>Bidder Username:</b> {winnerDetails?.Username || "Pending"}
        </h6>
        <h6>
          <b>Bidder Email:</b> {winnerDetails?.Email || "Pending"}
        </h6>
        <h6>
          <b>Bidder Phone:</b> {winnerDetails?.Phone || "Pending"}
        </h6>
        <h6>
          <b>Amount (€): </b>
          {state.WinningAmount || "Pending"}
        </h6>
      </div>
    </div>
  );
};
