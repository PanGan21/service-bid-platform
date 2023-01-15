import React, { useEffect, useState } from "react";
import { getMyRequests } from "../services/request";
import BootstrapTable, { ColumnDescription, PaginationOptions } from 'react-bootstrap-table-next';
import paginationFactory from "react-bootstrap-table2-paginator"



type Props = {}

// Not showing creatorId
const columns: ColumnDescription[] = [
    {
        dataField: "Id",
        text: "Id",
        sort: true
    },
    {
        dataField: "Title",
        text: "Title"
    },
    {
        dataField: "Postcode",
        text: "Postcode"
    },
    {
        dataField: "Info",
        text: "Info",
    },
    {
        dataField: "Deadline",
        text: "Deadline",
        sort: true
    },
    {
        dataField: "Status",
        text: "Status"
    }
]

export const Requests: React.FC<Props> = () => {
    const [requests, setRequests] = useState<Request[]>([]);

    useEffect(() => {
        getMyRequests().then(
            (response) => {
                if (response.data && response.data.length) {
                    setRequests(response.data);
                }

            },
            (error) => {
                // const _content = (error.response.data.error || error.message || JSON.stringify(error))
                // setContent(_content);
            }
        )
    }, [])
    return (
        <div>
            <BootstrapTable keyField="Id" data={requests} columns={columns} striped hover condensed pagination={paginationFactory({})} />
        </div>
    )

}