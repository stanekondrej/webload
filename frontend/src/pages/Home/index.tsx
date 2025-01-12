import { useState } from "preact/hooks";
import { useLocation } from "preact-iso";
import "./style.css";

export function Home() {
	const [targetId, setTargetId] = useState("");
	const location = useLocation();

	const keyHandler = (e: KeyboardEvent) => {
		if (e.key !== "Enter") {
			return;
		}

		e.preventDefault();
		location.route(targetId);
	};

	const inputHandler = (e: InputEvent) => {
		setTargetId(e.currentTarget["value"]);
	};

	return (
		<div class="home">
			<h1>Webload</h1>
			<p>
				<i>A dead simple system load tracker</i>
			</p>

			<div id="main-id-input">
				<textarea
					placeholder="Dashboard ID"
					onKeyDown={(e) => {
						keyHandler(e);
					}}
					onInput={(e) => {
						inputHandler(e);
					}}
				></textarea>
				<a href={window.location.href + targetId}>Go</a>
			</div>
		</div>
	);
}
