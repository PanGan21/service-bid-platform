import { useEffect, useState } from "react";
import { Column } from "react-table";
import { Pagination } from "../common/pagination";
import { AppTable } from "../common/table";
import { ROWS_PER_TABLE_PAGE } from "../constants";
import {
  countAuctionsByStatus,
  formatAuctions,
  getAuctionsByStatus,
  updateAuctionStatus,
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

const STATUS = "in progress";

const CLOSED_STATUS = "closed";

type Props = {};

export const InProgressAuctions: React.FC<Props> = () => {
  const [pageData, setPageData] = useState<{
    rowData: FormattedAuction[];
    isLoading: boolean;
    totalAuctions: number;
  }>({
    rowData: [],
    isLoading: false,
    totalAuctions: 0,
  });

  const [totalAuctions, setTotalAuctions] = useState<number>(0);
  const [currentPage, setCurrentPage] = useState<number>(1);

  useEffect(() => {
    setPageData((prevState) => ({
      ...prevState,
      rowData: [],
      isLoading: true,
    }));

    countAuctionsByStatus(STATUS).then((response) => {
      if (response.data && response.data) {
        setTotalAuctions(response.data);
      }
    });

    getAuctionsByStatus(STATUS, ROWS_PER_TABLE_PAGE, currentPage).then(
      (response) => {
        const auction = response.data || [];
        setPageData({
          isLoading: false,
          rowData: formatAuctions(auction),
          totalAuctions: totalAuctions,
        });
      }
    );
  }, [currentPage, totalAuctions]);

  const handleRowSelection = (auction: any) => {
    updateAuctionStatus(auction.Id, CLOSED_STATUS).then((response) => {
      window.location.reload();
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
        totalRows={pageData.totalAuctions}
        pageChangeHandler={setCurrentPage}
        rowsPerPage={ROWS_PER_TABLE_PAGE}
        currentPage={currentPage}
      />
    </div>
  );
};
