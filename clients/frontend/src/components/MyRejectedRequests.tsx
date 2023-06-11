import { useEffect, useState } from "react";
import { Column } from "react-table";
import { Pagination } from "../common/pagination";
import { AppTable } from "../common/table";
import { ROWS_PER_TABLE_PAGE } from "../constants";
import {
  countOwnRejectedRequests,
  getOwnRejectedRequests,
} from "../services/request";
import { Request } from "../types/request";

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
    Header: "Status",
    accessor: "Status",
  },
  {
    Header: "Rejection reason",
    accessor: "RejectionReason",
  },
];

type Props = {};

export const MyRejectedRequests: React.FC<Props> = () => {
  const [pageData, setPageData] = useState<{
    rowData: Request[];
    isLoading: boolean;
    totalRequest: number;
  }>({
    rowData: [],
    isLoading: false,
    totalRequest: 0,
  });

  const [totalRequest, setTotalRequest] = useState<number>(0);
  const [currentPage, setCurrentPage] = useState<number>(1);

  useEffect(() => {
    setPageData((prevState) => ({
      ...prevState,
      rowData: [],
      isLoading: true,
    }));

    countOwnRejectedRequests().then((response) => {
      if (response.data && response.data) {
        setTotalRequest(response.data);
      }
    });

    getOwnRejectedRequests(ROWS_PER_TABLE_PAGE, currentPage).then(
      (response) => {
        const requests = response.data || [];
        setPageData({
          isLoading: false,
          rowData: requests,
          totalRequest: totalRequest,
        });
      }
    );
  }, [currentPage, totalRequest]);

  const handleRowSelection = (request: any) => {};

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
        totalRows={pageData.totalRequest}
        pageChangeHandler={setCurrentPage}
        rowsPerPage={ROWS_PER_TABLE_PAGE}
        currentPage={currentPage}
      />
    </div>
  );
};
