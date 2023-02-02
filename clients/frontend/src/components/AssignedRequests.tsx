import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { Column } from "react-table";
import { Pagination } from "../common/pagination";
import { AppTable } from "../common/table";
import { ROWS_PER_TABLE_PAGE } from "../constants";
import {
  countRequestsByStatus,
  formatRequests,
  getRequestsByStatus,
} from "../services/request";
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

const STATUS = "assigned";

type Props = {};

export const AssignedRequests: React.FC<Props> = () => {
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
  const navigate = useNavigate();

  useEffect(() => {
    setPageData((prevState) => ({
      ...prevState,
      rowData: [],
      isLoading: true,
    }));

    countRequestsByStatus(STATUS).then((response) => {
      if (response.data && response.data) {
        setTotalRequests(response.data);
      }
    });

    getRequestsByStatus(STATUS, ROWS_PER_TABLE_PAGE, currentPage).then(
      (response) => {
        const requests = response.data || [];
        setPageData({
          isLoading: false,
          rowData: formatRequests(requests),
          totalRequests: totalRequests,
        });
      }
    );
  }, [currentPage, totalRequests]);

  const handleRowSelection = (request: any) => {
    navigate("/update-request-status", { state: request });
  };

  return (
    <div>
      <div style={{ height: "450px" }}>
        <AppTable
          columns={columns}
          data={pageData.rowData}
          isLoading={pageData.isLoading}
          onRowClick={(r) => handleRowSelection(r.values)}
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
