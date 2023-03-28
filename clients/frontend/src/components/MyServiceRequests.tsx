import { useState, useEffect } from "react";
import { AppTable, Column } from "../common/table";
import { Pagination } from "../common/pagination";
import { ROWS_PER_TABLE_PAGE } from "../constants";
import {
  countMyServiceRequests,
  formatAuctions,
  getMyServiceRequests,
} from "../services/auction";
import { FormattedAuction } from "../types/auction";

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

export const MyServiceRequests: React.FC<Props> = () => {
  const [pageData, setPageData] = useState<{
    rowData: FormattedAuction[];
    isLoading: boolean;
    totalServiceRequests: number;
  }>({
    rowData: [],
    isLoading: false,
    totalServiceRequests: 0,
  });
  const [totalServiceRequests, setTotalServiceRequests] = useState<number>(0);
  const [currentPage, setCurrentPage] = useState<number>(1);

  useEffect(() => {
    setPageData((prevState) => ({
      ...prevState,
      rowData: [],
      isLoading: true,
    }));

    countMyServiceRequests().then((response) => {
      if (response.data && response.data) {
        setTotalServiceRequests(response.data);
      }
    });

    getMyServiceRequests(ROWS_PER_TABLE_PAGE, currentPage).then((response) => {
      const auctions = response.data || [];
      setPageData({
        isLoading: false,
        rowData: formatAuctions(auctions),
        totalServiceRequests: totalServiceRequests,
      });
    });
  }, [currentPage, totalServiceRequests]);

  const handleRowSelection = () => {};

  return (
    <div>
      <div style={{ height: "450px" }}>
        <AppTable
          columns={columns}
          data={pageData.rowData}
          isLoading={pageData.isLoading}
          onRowClick={(r) => handleRowSelection()}
        />
      </div>
      <Pagination
        totalRows={pageData.totalServiceRequests}
        pageChangeHandler={setCurrentPage}
        rowsPerPage={ROWS_PER_TABLE_PAGE}
        currentPage={currentPage}
      />
    </div>
  );
};
