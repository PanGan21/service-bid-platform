import { Formik, Form, Field } from "formik";
import { useState } from "react";
import { NavigateFunction, useLocation, useNavigate } from "react-router-dom";
import request from "../assets/request.png";
import { updateRequestStatus } from "../services/request";
import { FormattedRequest } from "../types/request";

type Props = {};

const options = [
  { value: "closed", label: "closed" },
  { value: "in progress", label: "in progress" },
];

export const UpdateRequestStatus: React.FC<Props> = () => {
  const navigate: NavigateFunction = useNavigate();

  const [loading, setLoading] = useState<boolean>(false);
  const [message, setMessage] = useState<string>("");
  const { state }: { state: FormattedRequest } = useLocation();

  const initialValues: {
    status: string;
  } = {
    status: "",
  };

  const handleSubmit = async (formValue: { status: string }) => {
    const { status } = formValue;

    setMessage("");
    setLoading(true);

    try {
      await updateRequestStatus(state.Id, status);
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
        <img src={request} alt="profile-img" className="profile-img-card" />
        <Formik initialValues={initialValues} onSubmit={handleSubmit}>
          <Form>
            <div className="form-group">
              <div style={{ textAlign: "center" }}>
                Current status: {state.Status}
                <br />
                Request Id: {state.Id}
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
