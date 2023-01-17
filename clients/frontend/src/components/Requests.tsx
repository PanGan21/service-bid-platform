import { useState, useEffect } from "react";
import { AppTable, Column } from "../common/table";
import { Pagination } from "../common/pagination";
import { ROWS_PER_TABLE_PAGE } from "../constants";
import {
  countMyRequests,
  formatRequests,
  getMyRequests,
} from "../services/request";
import { PlusButton } from "../common/plus";
import { FormattedRequest } from "../types/request";

const columns: Column[] = [
  {
    Header: "Id",
    accessor: "Id",
  },
  {
    Header: "Title",
    accessor: "Title",
  },
  {
    Header: "Postcode",
    accessor: "Postcode",
  },
  {
    Header: "Info",
    accessor: "Info",
  },
  {
    Header: "Deadline",
    accessor: "Deadline",
  },
  {
    Header: "Status",
    accessor: "Status",
  },
];

type Props = {};

export const Requests: React.FC<Props> = () => {
  const [pageData, setPageData] = useState<{
    rowData: FormattedRequest[];
    isLoading: boolean;
    totalRequests: number;
  }>({
    rowData: [],
    isLoading: false,
    totalRequests: 0,
  });
  const [totalRequests, setTotalRequests] = useState<number>(0);
  const [currentPage, setCurrentPage] = useState<number>(1);

  useEffect(() => {
    setPageData((prevState) => ({
      ...prevState,
      rowData: [],
      isLoading: true,
    }));

    countMyRequests().then((response) => {
      if (response.data && response.data) {
        setTotalRequests(response.data);
      }
    });

    getMyRequests(ROWS_PER_TABLE_PAGE, currentPage).then((response) => {
      const requests = formatRequests(response.data) || [];
      setPageData({
        isLoading: false,
        rowData: requests,
        totalRequests: totalRequests,
      });
    });
  }, [currentPage, totalRequests]);

  return (
    <div>
      <div style={{ textAlign: "right" }}>
        <PlusButton navigation="/new-request" />
      </div>
      <div style={{ height: "450px" }}>
        <AppTable
          columns={columns}
          data={pageData.rowData}
          isLoading={pageData.isLoading}
        />
      </div>
      <Pagination
        totalRows={pageData.totalRequests}
        pageChangeHandler={setCurrentPage}
        rowsPerPage={ROWS_PER_TABLE_PAGE}
        currentPage={currentPage}
      />
    </div>
  );
};
