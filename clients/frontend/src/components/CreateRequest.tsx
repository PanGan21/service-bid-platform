import React, { useState } from "react";
import { NavigateFunction, useNavigate } from "react-router-dom";
import { Formik, Field, Form, ErrorMessage } from "formik";
import * as Yup from "yup";
import request from "../assets/request.png";

import { NewRequest } from "../types/request";
import { createRequest } from "../services/request";

type Props = {};

export const CreateRequest: React.FC<Props> = () => {
  const navigate: NavigateFunction = useNavigate();

  const [loading, setLoading] = useState<boolean>(false);
  const [message, setMessage] = useState<string>("");

  const today = new Date();
  const tomorrow = new Date(today);
  tomorrow.setDate(tomorrow.getDate() + 1);

  const initialValues: {
    title: string;
    postcode: string;
    info: string;
    deadline: string;
  } = {
    title: "",
    postcode: "",
    info: "",
    deadline: tomorrow.toLocaleDateString("en-US"),
  };

  const validationSchema = Yup.object().shape({
    title: Yup.string().required("This field is required!"),
    postcode: Yup.string()
      .matches(/^[0-9]{5}$/, "Invalid postcode")
      .required("This field is required!"),
    info: Yup.string().required("This field is required!"),
    deadline: Yup.date()
      .required("This field is required!")
      .test(
        "checkFutureDate",
        "Date should not be in the past!",
        (value: Date | undefined) => {
          if (!value) return false;
          return value > new Date();
        }
      ),
  });

  const handleSubmit = async (formValue: {
    title: string;
    postcode: string;
    info: string;
    deadline: string;
  }) => {
    const { title, postcode, info, deadline } = formValue;

    const request: NewRequest = {
      Title: title,
      Postcode: postcode,
      Info: info,
      Deadline: Date.parse(deadline),
    };

    setMessage("");
    setLoading(true);

    try {
      await createRequest(request);
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
        <img src={request} alt="profile-img" className="profile-img-card" />
        <Formik
          initialValues={initialValues}
          validationSchema={validationSchema}
          onSubmit={handleSubmit}
        >
          <Form>
            <div className="form-group">
              <label htmlFor="title">Title</label>
              <Field name="title" type="text" className="form-control" />
              <ErrorMessage
                name="title"
                component="div"
                className="alert alert-danger"
              />
            </div>

            <div className="form-group">
              <label htmlFor="postcode">Postcode</label>
              <Field name="postcode" type="postcode" className="form-control" />
              <ErrorMessage
                name="postcode"
                component="div"
                className="alert alert-danger"
              />
            </div>

            <div className="form-group">
              <label htmlFor="info">Info</label>
              <Field name="info" type="info" className="form-control" />
              <ErrorMessage
                name="info"
                component="div"
                className="alert alert-danger"
              />
            </div>

            <div className="form-group">
              <label htmlFor="deadline">Deadline</label>
              <Field name="deadline" type="deadline" className="form-control" />
              <ErrorMessage
                name="deadline"
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
