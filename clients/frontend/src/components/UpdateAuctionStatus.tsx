import { Formik, Form, Field } from "formik";
import { useEffect, useState } from "react";
import { NavigateFunction, useLocation, useNavigate } from "react-router-dom";
import auction from "../assets/auction.png";
import { rejectAuction, updateAuctionStatus } from "../services/auction";
import { getUserDetailsById } from "../services/auth";
import { FormattedAuction } from "../types/auction";
import { UserDetails } from "../types/user";

type Props = {};

type Option = {
  value: string;
  label: string;
};

const openStatusOptions: Option[] = [
  { value: "closed", label: "closed" },
  { value: "in progress", label: "in progress" },
];

const newStatusOptions: Option[] = [
  { value: "open", label: "open" },
  { value: "rejected", label: "reject" },
];

const allOptions: Option[] = [...openStatusOptions, ...newStatusOptions];

export const UpdateAuctionStatus: React.FC<Props> = () => {
  const navigate: NavigateFunction = useNavigate();
  const [options, setOptions] = useState<Option[]>(openStatusOptions);

  const initialValues: {
    status: string;
    rejectionReason: string;
  } = {
    status: "",
    rejectionReason: "",
  };

  const [loading, setLoading] = useState<boolean>(false);
  const [message, setMessage] = useState<string>("");
  const { state }: { state: FormattedAuction } = useLocation();
  const [winnerDetails, setWinnerDetails] = useState<UserDetails | undefined>(
    undefined
  );

  useEffect(() => {
    getUserDetailsById(state.WinnerId).then((response) => {
      if (response.data) {
        setWinnerDetails(response.data);
      }
    });
    if (state.Status === "new") {
      setOptions(newStatusOptions);
    } else if (state.Status === "assigned" || state.Status === "in progress") {
      setOptions(openStatusOptions);
    } else {
      setOptions(allOptions);
    }
  }, [state.Status, state.WinnerId]);

  const handleSubmit = async (formValue: {
    status: string;
    rejectionReason: string;
  }) => {
    const { status, rejectionReason } = formValue;

    setMessage("");
    setLoading(true);

    try {
      if (status === "rejected") {
        await rejectAuction(state.Id, rejectionReason);
      } else {
        await updateAuctionStatus(state.Id, status);
      }

      navigate("/new-auctions");
      window.location.reload();
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
                  <b>Current status:</b> {state.Status}
                  <br />
                  <b>Auction Id:</b> {state.Id}
                  <br />
                </div>
                {state.WinnerId && (
                  <div>
                    <b>Winner Username:</b> {winnerDetails?.Username}
                    <br />
                    <b>Winner Email:</b> {winnerDetails?.Email}
                    <br />
                    <b>Winner Phone:</b> {winnerDetails?.Phone}
                  </div>
                )}
                <br />
                <Field as="select" name="status" type="string">
                  <option value="">Select an option</option>
                  {options.map((option) => (
                    <option key={option.value} value={option.value}>
                      {option.label}
                    </option>
                  ))}
                </Field>
                {values.status === "rejected" && (
                  <div>
                    <br />
                    <Field
                      name="rejectionReason"
                      type="text"
                      placeholder="Enter reason for rejection"
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
