import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { Column } from "react-table";
import { Pagination } from "../common/pagination";
import { AppTable } from "../common/table";
import { ROWS_PER_TABLE_PAGE } from "../constants";
import {
  countOpenPastDeadlineAuctions,
  formatExtendedAuctions,
  getOpenPastDeadlineAuctions,
} from "../services/auction";
import { ExtendedFormattedAuction } from "../types/auction";

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

export const PendingAuctions: React.FC<Props> = () => {
  const [pageData, setPageData] = useState<{
    rowData: ExtendedFormattedAuction[];
    isLoading: boolean;
    totalAuctions: number;
  }>({
    rowData: [],
    isLoading: false,
    totalAuctions: 0,
  });
  const [totalAuctions, setTotalAuctions] = useState<number>(0);
  const [currentPage, setCurrentPage] = useState<number>(1);
  const navigate = useNavigate();

  useEffect(() => {
    setPageData((prevState) => ({
      ...prevState,
      rowData: [],
      isLoading: true,
    }));

    countOpenPastDeadlineAuctions().then((response) => {
      if (response.data) {
        setTotalAuctions(response.data);
      }
    });

    getOpenPastDeadlineAuctions(ROWS_PER_TABLE_PAGE, currentPage).then(
      (response) => {
        const auctions = response.data || [];
        setPageData({
          isLoading: false,
          rowData: formatExtendedAuctions(auctions),
          totalAuctions: totalAuctions,
        });
      }
    );
  }, [currentPage, totalAuctions]);

  const handleRowSelection = (auction: any) => {
    const fullAuction = pageData.rowData.find((r) => r.Id === auction.Id);
    navigate("/update-pending-auction", { state: fullAuction });
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
