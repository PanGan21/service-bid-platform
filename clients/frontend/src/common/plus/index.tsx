import { PlusCircle } from "react-bootstrap-icons";
import { useNavigate } from "react-router-dom";

interface PlusButtonProps {
  navigation: string;
}

export const PlusButton = ({ navigation }: PlusButtonProps) => {
  const navigate = useNavigate();

  return (
    <button className="btn btn-circle" onClick={() => navigate(navigation)}>
      <PlusCircle size={24} />
    </button>
  );
};
