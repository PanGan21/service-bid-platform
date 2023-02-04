import { MouseEventHandler } from "react";

export const ProfileImageBadge = ({
  src,
  badgeNumber,
  onBadgeClick,
}: {
  src: string;
  badgeNumber: number;
  onBadgeClick: MouseEventHandler;
}) => {
  return (
    <div style={{ position: "relative" }}>
      <img src={src} alt="" className="profile-img-card" />
      {badgeNumber >= 0 && (
        <div
          style={{
            position: "absolute",
            top: 0,
            right: 0,
            background: "transparent",
            color: "black",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            cursor: "pointer",
          }}
          onClick={onBadgeClick}
        >
          {badgeNumber}
        </div>
      )}
    </div>
  );
};
