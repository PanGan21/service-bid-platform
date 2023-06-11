import { useEffect, useState } from "react";
import { useLocation } from "react-router-dom";
import { Pagination } from "../common/pagination";
import { AppTable, Column } from "../common/table";
import { ROWS_PER_TABLE_PAGE } from "../constants";
import { countOwnAssignments, getOwnAssignments } from "../services/auction";
import { Auction } from "../types/auction";
import { User } from "../types/user";

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
    Header: "CreatorId",
    accessor: "CreatorId",
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
    Header: "WinningBidId",
    accessor: "WinningBidId",
  },
  {
    Header: "WinningAmount",
    accessor: "WinningAmount",
  },
];

type Props = {};

const handleRowSelection = (auction: any) => {};

export const Assignments: React.FC<Props> = () => {
  const [pageData, setPageData] = useState<{
    rowData: Auction[];
    isLoading: boolean;
    totalAssignments: number;
  }>({
    rowData: [],
    isLoading: false,
    totalAssignments: 0,
  });
  const [totalAssignments, setTotalAssignments] = useState<number>(0);
  const [currentPage, setCurrentPage] = useState<number>(1);
  const { state }: { state: User } = useLocation();

  useEffect(() => {
    setPageData((prevState) => ({
      ...prevState,
      rowData: [],
      isLoading: true,
    }));

    countOwnAssignments().then((response) => {
      if (response.data && response.data) {
        setTotalAssignments(response.data);
      }
    });

    getOwnAssignments(ROWS_PER_TABLE_PAGE, currentPage).then((response) => {
      const assignments = response.data || [];
      setPageData({
        isLoading: false,
        rowData: assignments,
        totalAssignments,
      });
    });
  }, [currentPage, totalAssignments]);

  return (
    <div>
      <div style={{ textAlign: "center" }}>
        <h3>
          The following is a list of auctions assigned to{" "}
          {state.Username !== "" ? state.Username : "you"}
        </h3>
      </div>
      <br />
      <div style={{ height: "450px" }}>
        <AppTable
          columns={columns}
          data={pageData.rowData}
          isLoading={pageData.isLoading}
          onRowClick={(r) => handleRowSelection(r.values)}
        />
      </div>
      <Pagination
        totalRows={pageData.totalAssignments}
        pageChangeHandler={setCurrentPage}
        rowsPerPage={ROWS_PER_TABLE_PAGE}
        currentPage={currentPage}
      />
    </div>
  );
};
