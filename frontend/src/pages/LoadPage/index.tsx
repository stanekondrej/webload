import { Signal, signal } from "@preact/signals";
import { useRoute } from "preact-iso";
import "./style.css";
import { useEffect } from "preact/hooks";
import { Stats } from "../../../proto/messages";

const url = `${import.meta.env["VITE_WEBLOAD_SERVER"]}/query`;

type numSignal = Signal<number>;

export function LoadPage() {
	const route = useRoute();

	const cpu: numSignal = signal();
	const mem: numSignal = signal();
	const uptime: numSignal = signal();

	useEffect(() => {
		const conn = new WebSocket(url);

		conn.onerror = (e) => {
			console.error(e);
		};

		conn.onopen = () => {
			conn.send(id);
		};

		conn.onclose = (e) => {
			console.log("Connection closed", e);
		};

		conn.onmessage = async (e) => {
			const data = new Blob(e.data);

			let message: Stats;
			try {
				message = Stats.decode(await data.bytes());
			} catch (err) {
				console.error(err);
				return;
			}

			cpu.value = message.cpuUsage;
		};

		return () => {
			conn.close();
		};
	});
	const id = route.path.slice(1);

	return (
		<div id="info-container">
			<h1>Viewing ID {id}</h1>
			<table>
				<thead>
					<tr>
						<th>CPU</th>
						<th>MEM</th>
						<th>Uptime</th>
					</tr>
				</thead>

				<tbody>
					<tr>
						<td>{cpu}</td>
						<td>{mem}</td>
						<td>{uptime}</td>
					</tr>
				</tbody>
			</table>
		</div>
	);
}
