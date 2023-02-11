import { useLocation } from "react-router-dom";
import bid from "../assets/bid.png";
import { Bid } from "../types/bid";

type Props = {};

export const AssignAuction: React.FC<Props> = () => {
  const { state }: { state: Bid } = useLocation();

  return (
    <div className="col-md-12">
      <div className="card card-container">
        <img src={bid} alt="profile-img" className="profile-img-card" />
        <h3 style={{ textAlign: "center" }}>Winning bid</h3>
        <h6>Bidder Id: {state.CreatorId}</h6>
        <h6>Auction Id: {state.AuctionId}</h6>
        <h6>Amount: {state.Amount} €</h6>
      </div>
    </div>
  );
};
