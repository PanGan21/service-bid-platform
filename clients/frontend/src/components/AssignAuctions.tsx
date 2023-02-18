import { useLocation } from "react-router-dom";
import bid from "../assets/bid.png";
import { Auction } from "../types/auction";

type Props = {};

export const AssignAuction: React.FC<Props> = () => {
  const { state }: { state: Auction } = useLocation();

  return (
    <div className="col-md-12">
      <div className="card card-container">
        <img src={bid} alt="profile-img" className="profile-img-card" />
        <h3 style={{ textAlign: "center" }}>Winning bid</h3>
        <h6>Bidder Id: {state.WinnerId}</h6>
        <h6>Bid Id: {state.WinningBidId}</h6>
        <h6>Auction Id: {state.Id}</h6>
        <h6>Amount: {state.WinningAmount} €</h6>
      </div>
    </div>
  );
};
