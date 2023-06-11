import { Formik, Form, Field } from "formik";
import { useState } from "react";
import { NavigateFunction, useLocation, useNavigate } from "react-router-dom";
import auction from "../assets/auction.png";
import { updateAuctionDeadline, updateWinner } from "../services/auction";
import { FormattedAuction } from "../types/auction";

type Props = {};

type Option = {
  value: string;
  label: string;
};

const options: Option[] = [
  { value: "assign", label: "assign winner" },
  { value: "extend", label: "extend deadline" },
];

export const UpdatePendingAuction: React.FC<Props> = () => {
  const navigate: NavigateFunction = useNavigate();

  const initialValues: {
    action: string;
    days: number;
  } = {
    action: "",
    days: 0,
  };

  const [loading, setLoading] = useState<boolean>(false);
  const [message, setMessage] = useState<string>("");
  const { state }: { state: FormattedAuction } = useLocation();

  const handleSubmit = async (formValue: { action: string; days: number }) => {
    const { action, days } = formValue;

    setMessage("");
    setLoading(true);

    try {
      if (action === "assign") {
        updateWinner(state.Id)
          .then((response) => {
            if (response.data && response.data) {
              navigate("/assigned-auction", { state: response.data });
            }
          })
          .catch((error) => {
            if (error.response.data.error === "Could not find winning bid") {
              alert("Bids not found for this auction!");
              window.location.reload();
            }
          });
      } else {
        updateAuctionDeadline(state.Id, days).then((response) => {
          if (response.data && response.data) {
            navigate("/open-auctions", { state: response.data });
          }
        });
      }
    } catch (error: any) {
      const resMessage =
        (error.response &&
          error.response.data &&
          error.response.data.message) ||
        error.message ||
        error.toString();

      setLoading(false);
      setMessage(resMessage);
    }
  };

  return (
    <div className="col-md-12">
      <div className="card card-container">
        <img src={auction} alt="profile-img" className="profile-img-card" />
        <Formik initialValues={initialValues} onSubmit={handleSubmit}>
          {({ values }) => (
            <Form>
              <div className="form-group">
                <div style={{ textAlign: "left" }}>
                  <b>Auction Id:</b> {state.Id}
                  <br />
                  <b>Current status:</b> {state.Status}
                  <br />
                  <b>Deadline:</b> {state.Deadline}
                  <br />
                </div>
                <br />
                <Field as="select" name="action" type="string">
                  <option value="">Select an option</option>
                  {options.map((option) => (
                    <option key={option.value} value={option.value}>
                      {option.label}
                    </option>
                  ))}
                </Field>
                {values.action === "extend" && (
                  <div>
                    <br />
                    Enter number of days to extend Auction
                    <Field
                      name="days"
                      type="number"
                      min="0"
                      max="30"
                      style={{ width: "100%" }}
                    />
                  </div>
                )}
                <br />
                <br />
                <div className="form-group">
                  <button
                    type="submit"
                    className="btn btn-primary btn-block"
                    disabled={loading}
                  >
                    {loading && (
                      <span className="spinner-border spinner-border-sm"></span>
                    )}
                    <span>Submit</span>
                  </button>
                </div>
                {message && (
                  <div className="form-group">
                    <div className="alert alert-danger" role="alert">
                      {message}
                    </div>
                  </div>
                )}
              </div>
            </Form>
          )}
        </Formik>
      </div>
    </div>
  );
};
