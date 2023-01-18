import { ReactNode, useState } from "react";

type Props = {
  children: ReactNode;
  onClick: React.MouseEventHandler<HTMLTableRowElement>;
};

export const TableRow = ({ children, onClick }: Props) => {
  const [opacity, setOpacity] = useState(1);

  return (
    <tr
      style={{ opacity, cursor: "pointer" }}
      onMouseEnter={() => setOpacity(0.5)}
      onMouseLeave={() => setOpacity(1)}
      onClick={onClick}
    >
      {children}
    </tr>
  );
};
