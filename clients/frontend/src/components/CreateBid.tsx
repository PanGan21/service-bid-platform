import React, { useState } from "react";
import { NavigateFunction, useLocation, useNavigate } from "react-router-dom";
import { FormattedRequest } from "../types/request";
import * as Yup from "yup";
import { NewBid } from "../types/bid";
import { Formik, Field, Form, ErrorMessage } from "formik";
import bid from "../assets/bid.png";
import { createBid } from "../services/bid";

type Props = {};

export const CreateBid: React.FC<Props> = () => {
  const navigate: NavigateFunction = useNavigate();
  const { state }: { state: FormattedRequest } = useLocation();

  const [loading, setLoading] = useState<boolean>(false);
  const [message, setMessage] = useState<string>("");

  const initialValues: {
    amount: number;
    requstId: string;
  } = {
    amount: 0,
    requstId: "",
  };

  const validationSchema = Yup.object().shape({
    amount: Yup.number().moreThan(0, "The number must be greater than 0"),
  });

  const handleSubmit = async (formValue: { amount: number }) => {
    const { amount } = formValue;

    const newBid: NewBid = {
      Amount: amount,
      RequestId: state.Id,
    };

    setMessage("");
    setLoading(true);

    try {
      await createBid(newBid);
      navigate("/all-requests");
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
        <img src={bid} alt="profile-img" className="profile-img-card" />
        <Formik
          initialValues={initialValues}
          validationSchema={validationSchema}
          onSubmit={handleSubmit}
        >
          <Form>
            <div className="form-group">
              <label htmlFor="amount">
                Amount in € for request in {state.Postcode}
              </label>
              <Field name="amount" type="number" className="form-control" />
              <ErrorMessage
                name="amount"
                component="div"
                className="alert alert-danger"
              />
            </div>

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
          </Form>
        </Formik>
      </div>
    </div>
  );
};
