import { Formik, Form, Field } from "formik";
import { useEffect, useState } from "react";
import { NavigateFunction, useLocation, useNavigate } from "react-router-dom";
import auction from "../assets/auction.png";
import { updateAuctionStatus } from "../services/auction";
import { FormattedAuction } from "../types/auction";

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

  const [loading, setLoading] = useState<boolean>(false);
  const [message, setMessage] = useState<string>("");
  const { state }: { state: FormattedAuction } = useLocation();

  const initialValues: {
    status: string;
  } = {
    status: "",
  };

  useEffect(() => {
    if (state.Status === "new") {
      setOptions(newStatusOptions);
    } else if (state.Status === "assigned" || state.Status === "in progress") {
      setOptions(openStatusOptions);
    } else {
      setOptions(allOptions);
    }
  }, [state.Status]);

  const handleSubmit = async (formValue: { status: string }) => {
    const { status } = formValue;

    setMessage("");
    setLoading(true);

    try {
      await updateAuctionStatus(state.Id, status);
      navigate("/admin");
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
          <Form>
            <div className="form-group">
              <div style={{ textAlign: "center" }}>
                Current status: {state.Status}
                <br />
                Auction Id: {state.Id}
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
        </Formik>
      </div>
    </div>
  );
};
