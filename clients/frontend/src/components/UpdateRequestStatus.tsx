import { Formik, Form, Field } from "formik";
import { useState } from "react";
import { NavigateFunction, useLocation, useNavigate } from "react-router-dom";
import auction from "../assets/auction.png";
import { approveRequest, rejectRequest } from "../services/request";
import { Request } from "../types/request";

type Props = {};

type Option = {
  value: string;
  label: string;
};

const options: Option[] = [
  { value: "approved", label: "approved" },
  { value: "rejected", label: "rejected" },
];

export const UpdateRequestStatus: React.FC<Props> = () => {
  const navigate: NavigateFunction = useNavigate();

  const initialValues: {
    status: string;
    rejectionReason: string;
    days: number;
  } = {
    status: "",
    rejectionReason: "",
    days: 0,
  };

  const [loading, setLoading] = useState<boolean>(false);
  const [message, setMessage] = useState<string>("");
  const { state }: { state: Request } = useLocation();

  const handleSubmit = async (formValue: {
    status: string;
    rejectionReason: string;
    days: number;
  }) => {
    const { status, rejectionReason, days } = formValue;

    setMessage("");
    setLoading(true);

    try {
      if (status === "rejected") {
        await rejectRequest(state.Id, rejectionReason);
      } else {
        await approveRequest(state.Id, days);
      }

      navigate("/new-service-requests");
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
                  <b>Request Id:</b> {state.Id}
                  <br />
                </div>
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
                    Enter reason for rejection
                    <Field
                      name="rejectionReason"
                      type="text"
                      style={{ width: "100%" }}
                    />
                  </div>
                )}
                {values.status === "approved" && (
                  <div>
                    <br />
                    Enter number of days to keep Auction open
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
