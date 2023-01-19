import { useState, useEffect } from "react";
import { AppTable, Column } from "../common/table";
import { Pagination } from "../common/pagination";
import { ROWS_PER_TABLE_PAGE } from "../constants";
import {
  countAllRequests,
  formatRequests,
  getAllRequests,
} from "../services/request";
import { PlusButton } from "../common/plus";
import { FormattedRequest } from "../types/request";
import { useNavigate } from "react-router-dom";

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

export const AllRequests: React.FC<Props> = () => {
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

    countAllRequests().then((response) => {
      if (response.data && response.data) {
        setTotalRequests(response.data);
      }
    });

    getAllRequests(ROWS_PER_TABLE_PAGE, currentPage).then((response) => {
      const requests = response.data || [];
      setPageData({
        isLoading: false,
        rowData: formatRequests(requests),
        totalRequests: totalRequests,
      });
    });
  }, [currentPage, totalRequests]);

  const handleRowSelection = (request: any) => {
    navigate("/new-bid", { state: request });
  };

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
