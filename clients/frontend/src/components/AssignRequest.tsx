import { useLocation } from "react-router-dom";
import bid from "../assets/bid.png";
import { Bid } from "../types/bid";

type Props = {};

export const AssignRequest: React.FC<Props> = () => {
  const { state }: { state: Bid } = useLocation();

  return (
    <div className="col-md-12">
      <div className="card card-container"></div>
      <img src={bid} alt="profile-img" className="profile-img-card" />
      {JSON.stringify(state)}
    </div>
  );
};
