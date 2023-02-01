import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { Column } from "react-table";
import { Pagination } from "../common/pagination";
import { AppTable } from "../common/table";
import { ROWS_PER_TABLE_PAGE } from "../constants";
import {
  countOpenPastDeadlineRequests,
  formatExtendedRequests,
  getOpenPastDeadlineRequests,
  updateWinner,
} from "../services/request";
import { ExtendedFormattedRequest } from "../types/request";

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
  {
    Header: "# Bids",
    accessor: "BidsCount",
  },
];

type Props = {};

export const PendingRequests: React.FC<Props> = () => {
  const [pageData, setPageData] = useState<{
    rowData: ExtendedFormattedRequest[];
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

    countOpenPastDeadlineRequests().then((response) => {
      if (response.data && response.data) {
        setTotalRequests(response.data);
      }
    });

    getOpenPastDeadlineRequests(ROWS_PER_TABLE_PAGE, currentPage).then(
      (response) => {
        const requests = response.data || [];
        setPageData({
          isLoading: false,
          rowData: formatExtendedRequests(requests),
          totalRequests: totalRequests,
        });
      }
    );
  }, [currentPage, totalRequests]);

  const handleRowSelection = (request: any) => {
    updateWinner(request.Id)
      .then((response) => {
        if (response.data && response.data) {
          navigate("/assign-request", { state: response.data });
        }
      })
      .catch((error) => {
        console.log(error.response);
        if (error.response.data.error === "Could not find winning bid") {
          alert("Bids not found for this request!");
        }
      });
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
